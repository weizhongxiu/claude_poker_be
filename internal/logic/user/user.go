package user

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"

	"claude-test/internal/model/entity"
	"claude-test/utility/jwt"
)

func Register(ctx context.Context, phone, password, nickname string) (userID int64, token string, err error) {
	// 检查手机号是否已注册
	count, err := g.DB().Model("users").Where("phone", phone).Count()
	if err != nil {
		return
	}
	if count > 0 {
		err = gerror.New("手机号已注册")
		return
	}

	uid := fmt.Sprintf("%d", time.Now().UnixMicro())[4:14]
	hash, err := gmd5.EncryptString(password + uid)
	if err != nil {
		return
	}

	result, err := g.DB().Model("users").Data(g.Map{
		"uid":      uid,
		"nickname": nickname,
		"phone":    phone,
		"password": hash,
		"status":   1,
		"avatar":   fmt.Sprintf("https://api.dicebear.com/7.x/adventurer/svg?seed=%s", uid),
	}).Insert()
	if err != nil {
		return
	}
	userID, err = result.LastInsertId()
	if err != nil {
		return
	}

	// 初始化钱包，赠送 10000 筹码
	_, err = g.DB().Model("user_wallets").Insert(g.Map{
		"user_id": userID,
		"chips":   10000,
	})
	if err != nil {
		return
	}

	token, err = jwt.Generate(userID)
	return
}

func Login(ctx context.Context, phone, password string) (user *entity.User, wallet *entity.UserWallet, token string, err error) {
	user = &entity.User{}
	err = g.DB().Model("users").Where("phone", phone).Where("status", 1).Scan(user)
	if err != nil {
		return
	}
	if user.Id == 0 {
		err = gerror.New("账号不存在")
		return
	}

	hash, _ := gmd5.EncryptString(password + user.Uid)
	if hash != user.Password {
		err = gerror.New("密码错误")
		return
	}

	// 更新最后登录时间
	now := time.Now()
	_, _ = g.DB().Model("users").Where("id", user.Id).Update(g.Map{"last_login_at": now})

	wallet = &entity.UserWallet{}
	_ = g.DB().Model("user_wallets").Where("user_id", user.Id).Scan(wallet)

	token, err = jwt.Generate(user.Id)

	// Token 存 Redis，支持踢下线
	_, _ = g.Redis().Set(ctx, fmt.Sprintf("token:%d", user.Id), token)
	_, _ = g.Redis().Expire(ctx, fmt.Sprintf("token:%d", user.Id), 7*24*3600)

	_ = grand.S(1) // suppress unused import
	return
}

func Profile(ctx context.Context, userID int64) (user *entity.User, wallet *entity.UserWallet, err error) {
	user = &entity.User{}
	err = g.DB().Model("users").Where("id", userID).Scan(user)
	if err != nil || user.Id == 0 {
		err = gerror.New("用户不存在")
		return
	}
	wallet = &entity.UserWallet{}
	_ = g.DB().Model("user_wallets").Where("user_id", userID).Scan(wallet)
	return
}
