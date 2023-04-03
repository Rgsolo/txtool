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
