package sender

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"reflect"
	"testing"
)

func Test_getTransactionSender(t *testing.T) {
	decode, _ := hexutil.Decode("0x02f8768189827af9850826299e0085654e543f21825208942c03058178e4f93209a78e8fd9b6fd44ccd77526877db75624aefbf080c001a04ef06811b9ca01961ab8ca34f318de479ed2d4753d9666889c970f492c86345ea077458523e8c0099e02100682596d94b1b604fa8df0e4f618b4573488b5fc4359")
	transaction := new(types.Transaction)
	_ = transaction.UnmarshalBinary(decode)
	type args struct {
		tx *types.Transaction
	}
	tests := []struct {
		name    string
		args    args
		want    common.Address
		wantErr bool
	}{
		{
			name: "getTransactionSender",
			args: args{tx: transaction},
			want: common.HexToAddress("0x5c7beD3Cca42e4562877eD88B9Aa0F5898Ed59B0"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTransactionSender(tt.args.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTransactionSender() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTransactionSender() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//
//func Sum(x int, y int) int {
//	if x < -1 {
//		panic("error x")
//	}
//	return x + y
//}
//
//func FuzzSum(f *testing.F) {
//	f.Add(-2, 3)
//	f.Fuzz(func(t *testing.T, x int, y int) {
//		z := Sum(x, y)
//		if z != x+y {
//			t.Error()
//		}
//	})
//}
//
//func TestSum(t *testing.T) {
//	assertions := assert.New(t)
//
//	result := Sum(2, 3)
//	assertions.Equal(5, result, "they should be equal")
//
//	result = Sum(-1, 1)
//	assertions.Equal(0, result, "they should be equal")
//
//	result = Sum(-1, -1)
//	assertions.Equal(-2, result, "they should be equal")
//}
//
//type User struct {
//	ID   int
//	Name string
//}
//
//type UserFetcher interface {
//	FetchUser(id int) (User, error)
//}
//
//func GetUserName(id int, uf UserFetcher) (string, error) {
//	user, err := uf.FetchUser(id)
//	if err != nil {
//		return "", err
//	}
//	return user.Name, nil
//}
//
//type MockUserFetcher struct {
//	mock.Mock
//}
//
//func (m *MockUserFetcher) FetchUser(id int) (User, error) {
//	args := m.Called(id)
//	return args.Get(0).(User), args.Error(1)
//}
//
//func TestGetUserName(t *testing.T) {
//	mockUserFetcher := new(MockUserFetcher)
//
//	user := User{
//		ID:   1,
//		Name: "John Doe",
//	}
//
//	// 设置当 FetchUser 被调用时返回什么
//	mockUserFetcher.On("FetchUser", 1).Return(user, nil)
//
//	name, err := GetUserName(1, mockUserFetcher)
//
//	assert.NoError(t, err)
//	assert.Equal(t, "John Doe", name)
//
//	// 确保 FetchUser 方法被调用了
//	mockUserFetcher.AssertCalled(t, "FetchUser", 1)
//}
