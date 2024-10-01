package sender

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
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

func TestEncryptAES(t *testing.T) {
	plaintext := "308204bd020100300d06092a864886f70d0101010500048204a7308204a30201000282010100b9d1910eee184e5cc0648a7d33d0e33812a802cd0dd3db2e25136c94e83c587a47ed0727075f3d228fcb1ef95fa2b0a3f83c41566633e6958caaca934ac87123fce1ba0737941951cf5d3447568a12fc18b898765e738c160032b37c29b0e5c0aeb81383161145879978b722e49f889ac26c761467e458b2a2801176073a8922d86f1f55e903a3e37ed8aa43707e1827442904572b31057233a542b05f66610e14d52dda64f8714612777556bd4c8484c4073e3b33364d0d1bbe56ca6239573b8b79d3e1b6fadb5632196b6e9fac6c304ba8897530a067f6f96bfb35ade22e5757c62aa3933ed7b5efb0190e5bcb9c0fb8e5f75c061eee45694001679d75e8ef02030100010282010011c421b465bb5932b10a6bdb50aaf62e944a100a7ef9f488d2eeaa810a3b4ad25632296ee7db8942d6b0bb0368cb6b4c221dd0b96c082651c2234a3f0ef55f2bebafc1539352cba0f0cc9e84fb9733fd7a779bcc2577b2bb1fc5b93773dcfffa8e39f7539f36838955791f396cd67bff1ffbb2c67cb06e7295eadcc0862c68940c95e6a36ff0168f0685ae9db21b6f588b0dd710d9f3751a8655dc5b7ec715b4f4929e418a0f29da29c0bf66f54eb293706c4ccfa055b62727cf567868219020f50f53d8c7d7e01f03d08f361992e6bb3097d72f6c48f3146f3fac24c5fdbab512031c86c8100b0ad3c16ecfaeeed4c62171ee7f95b6be3201a4a875a74fe90102818100d2812ea34f4cd103cffcdc6eac0dc6e81522fd07f3ad75f0edf08990e519aaee2398c6161340a8e560f0b5973a7dc016b1e9553483d86aa5b743684565b350caddeb04960d10d27c35c89d480e555ec7019c859a785b24a1a2b980a96d6707023130c25d7343d46b96b32c5ab4d51c8c99704eb33bce74e2150d2d64d474cac102818100e1fa8e2e73545786d1b860721ce2a7351dee2a021f3bb194408adbe059c14b00104db1b56f59a7cd777976f4acfc0a349655e9d8a290b0479bc73cf8e8d51c5140a02d571dabef91053d50f15787937eb99f2f6fb3455ee522e35dd961a297683ddd209e2f63e268389eb05a84220fdb1dd2a5e6a2ff80a130ca65c52f3e0faf028180187f7639054290549c40f63a5f059f6f64fe546a377ea96a2796c5bbba1dd999ab44dd50b65ee10908d61a9c05ef9a8a499c39114a82c62e90fc64472745ce123def5af24784fca9fdec61f97fc989a52957d8e898372b353065dc465b781105bd49ca64ebc42a15774d54cd1d9c6b9d25423fb6763059c3f1e53db22d53864102818100c95322e6c942e3f3add618beb173504a674eeec813316864cf17a70a7a8c559849ee1e1ba9877392ff150ce0b1589e72f958b34c9890552c86e17b35baa15a681d2e57169ccdb852381bc7882c561216bda2cf6d3186e52338f0dd84b4925957ba7dcdf87ff9fa14a614e2c2d1a3530300cfd193e7b5bfd05b9a860ea5a721eb02818040f92420303c754f2cbe4339be9362721e85cf68d23bdc47ecaa2ec58905e4c8a28dc5d48038adeb4c978410ec0d5e083b9f1d216cc93aa5f39bbf4fd80c84e3e9c4d96e4315c5d19532efe9eb9bb66de07a327c5ec27dc5b6679c15706cd240050d317255dc55b28f5ca68b740bcd96966327d2bf3a11a8b4f561af9ea6c34f"
	key := []byte("12345678901234567890123456789012")
	ciphertext := encryptAES(key, plaintext)
	t.Log(ciphertext)
	decryptedText := decryptAES(key, ciphertext)
	if decryptedText != plaintext {
		t.Errorf("decrypted text does not match plaintext")
	}
}

func decryptAES(key []byte, ct string) string {
	ciphertext, _ := hex.DecodeString(ct)
	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	size := c.BlockSize()
	pt := make([]byte, len(ciphertext))
	for bs, be := 0, size; bs < len(ciphertext); bs, be = bs+size, be+size {
		c.Decrypt(pt[bs:be], ciphertext[bs:be])
	}
	pt = unPaddingPKCS7(pt)
	return string(pt[:])
}
func encryptAES(key []byte, plaintext string) string {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	pt := []byte(plaintext)
	data := paddingPKCS7(pt, block.BlockSize())
	ciphertext := make([]byte, len(data))
	size := block.BlockSize()
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(ciphertext[bs:be], data[bs:be])
	}
	// fmt.Println(hex.EncodeToString(ciphertext))
	return hex.EncodeToString(ciphertext)
}
func paddingPKCS7(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}
func unPaddingPKCS7(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
