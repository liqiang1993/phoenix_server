package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	pb "github.com/lucky-cheerful-man/phoenix_apis/protobuf3.pb/user_info_manage"
	"github.com/lucky-cheerful-man/phoenix_server/src/config"
	"github.com/lucky-cheerful-man/phoenix_server/src/constant"
	"github.com/lucky-cheerful-man/phoenix_server/src/gmysql"
	"github.com/lucky-cheerful-man/phoenix_server/src/gredis"
	"github.com/lucky-cheerful-man/phoenix_server/src/log"
	"github.com/lucky-cheerful-man/phoenix_server/src/util"
)

type UserService struct {
}

type UserInfo struct {
	Nickname  string
	ImagePath string
}

// Register 注册接口
func (u *UserService) Register(_ context.Context, in *pb.RegisterRequest, _ *pb.RegisterResponse) error {
	password := util.EncodeMD5(in.Password + config.ReferGlobalConfig().AppSetting.Salt)
	err := gmysql.DBOperate.InsertUser(in.Name, password)
	if err != nil {
		log.Warn("%s InsertUser failed, err:%s", in.RequestID, err)
		return errors.New(constant.InsertDBError.Msg)
	}

	// 存储缓存时出现错误，不影响正常流程
	u.setCacheInfo(in.RequestID, in.Name, in.Name, "")

	return nil
}

// Auth 登陆认证接口
func (u *UserService) Auth(_ context.Context, in *pb.AuthRequest, out *pb.AuthResponse) error {
	password := util.EncodeMD5(in.Password + config.ReferGlobalConfig().AppSetting.Salt)
	res, nickname, imagePath, err := gmysql.DBOperate.CheckAuth(in.Name, password)
	if err != nil || !res {
		log.Warn("%s CheckAuth failed, err:%s, res:%v", in.RequestID, err, res)
		return errors.New("auth failed")
	}
	out.Image = imagePath
	out.Nickname = nickname
	return nil
}

// GetProfile 查询用户属性信息
func (u *UserService) GetProfile(_ context.Context, in *pb.GetProfileRequest, out *pb.GetProfileResponse) error {
	// 优先从缓存查询
	userInfoPtr := new(UserInfo)
	res, err := gredis.CacheOperate.Get(in.Name)
	if err == nil {
		err = json.Unmarshal(res, userInfoPtr)
		if err == nil {
			out.Nickname = userInfoPtr.Nickname
			out.ImageID = userInfoPtr.ImagePath
			return nil
		}
	}

	log.Info("%s get user info from redis failed:%s, res:%s", in.RequestID, err, res)

	// 缓存中不存在数据时，从db查询
	nickname, imagePath, err := gmysql.DBOperate.GetProfile(in.Name)
	if err != nil {
		log.Warn("%s GetProfile failed, err:%s,", in.RequestID, err)
		return err
	}

	// db查询成功后，更新缓存
	u.setCacheInfo(in.RequestID, in.Name, nickname, imagePath)
	out.Nickname = nickname
	out.ImageID = imagePath
	return nil
}

// GetHeadImage 查询头像图片
func (u *UserService) GetHeadImage(_ context.Context, in *pb.GetHeadImageRequest, out *pb.GetHeadImageResponse) error {
	file, err := os.Open(config.ReferGlobalConfig().AppSetting.RootPictureDir + in.ImageID)
	if err != nil {
		log.Warn("%s open file failed, err:%s", in.RequestID, err)
		return err
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			log.Warn("%s file close failed:%s", in.RequestID, closeErr)
		}
	}()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Warn("%s ReadAll file failed, err:%s", in.RequestID, err)
		return err
	}

	out.Image = content
	return nil
}

// EditProfile 编辑用户属性信息
func (u *UserService) EditProfile(_ context.Context, in *pb.EditProfileRequest, _ *pb.EditProfileResponse) error {
	var err error
	var path string
	var imageID string

	if len(in.Image) != 0 {
		// todo 更新图片成功后，需要删除原有的图片
		imageID = fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Int()) //nolint:gosec
		path = config.ReferGlobalConfig().AppSetting.RootPictureDir + imageID
		err = ioutil.WriteFile(path, in.Image, 0644) //nolint:gomnd,gosec
		if err != nil {
			log.Warn("%s write file failed, err:%s", in.RequestID, err)
			return err
		}
	}

	err = gmysql.DBOperate.EditProfile(in.Name, imageID, in.Nickname)
	if err != nil {
		log.Warn("%s EditProfile failed, err:%s", in.RequestID, err)
		// 更新失败时删除刚存储的图片
		fileErr := os.Remove(path)
		if fileErr != nil {
			log.Warn("%s os.Remove failed, fileErr:%s", in.RequestID, fileErr)
		}
		return err
	}

	// 删除缓存中的内容
	_, err = gredis.CacheOperate.Delete(in.Name)
	if err != nil {
		log.Warn("%s delete cache failed:%s", in.RequestID, err)
	}

	return nil
}

func (u *UserService) setCacheInfo(requestID string, name string, nickname string, imagePath string) {
	userInfo := UserInfo{
		Nickname:  nickname,
		ImagePath: imagePath,
	}
	res, err := json.Marshal(userInfo)
	if err == nil {
		err = gredis.CacheOperate.Set(name, res, config.ReferGlobalConfig().RedisSetting.ExpireTimeSecond)
		if err != nil {
			log.Warn("%s set cache failed:%s, name:%s, res:%s", requestID, err, name, string(res))
		}
	} else {
		log.Warn("%s json.Marshal failed:%s, userInfo:%+v", requestID, err, userInfo)
	}
}
