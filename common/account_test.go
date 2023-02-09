package common_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana/common"
	"github.com/stretchr/testify/require"
)

func TestSLIP10Compatibility_12words(t *testing.T) {
	mnemonic := "response photo senior language wave property trip purse bench arena casual noodle"
	account, err := common.DeriveAccountFromMnemonicBip44(mnemonic)
	require.NoError(t, err)

	{
		expectedPrivateKeyHex := "fc680ab5ec8b348ac8a6169327b3d2968dd037b43c93575338692088eb0c21626fe66b38dd599fdc817fe5cc3b9e8e524a826ba031fe99d409db130b79b825db"
		actualPrivateKeyHex := hex.EncodeToString(account.PrivateKey)
		require.Equal(t, expectedPrivateKeyHex, actualPrivateKeyHex)
	}

	{
		expectedAddr := "8Xp3CxmnwTbjYNKwsKEqgCSozqGWcDZHCWtAnxWb86oc"
		actualAddr := account.PublicKey.ToBase58()
		require.Equal(t, expectedAddr, actualAddr)
	}
}

func TestSLIP10Compatibility_24words(t *testing.T) {
	mnemonic := "diagram another jealous will cost ship goose blind elevator anxiety crazy cheese " +
		"cherry jeans rhythm february fat broom tattoo artwork cluster damp maple scorpion"
	account, err := common.DeriveAccountFromMnemonicBip44(mnemonic)
	require.NoError(t, err)

	{
		expectedPrivateKeyHex := "623c0c7fbdd49b93a33aef2a1eada0f1f9ee7d06f958194ed8a7a1fa6b76d47f" +
			"7541f1271fecbb9fad2501077b20779d0fc5448c45fcd549ac7c2ba81cf676b0"
		actualPrivateKeyHex := hex.EncodeToString(account.PrivateKey)
		require.Equal(t, expectedPrivateKeyHex, actualPrivateKeyHex)
	}

	{
		expectedAddr := "8tj2AYrV3bNHaayZuTiQs5vShJH57PtnBsDYJT7QBEK9"
		actualAddr := account.PublicKey.ToBase58()
		require.Equal(t, expectedAddr, actualAddr)
	}
}

func TestAccountFromMnemonicBip39_12Words(t *testing.T) {
	mnemonic, err := common.NewMnemonic(common.MnemonicLength12)
	require.NoError(t, err)
	require.NotEmpty(t, mnemonic)

	account, err := common.DeriveAccountFromMnemonicBip39(mnemonic)
	require.NoError(t, err)
	require.NotNil(t, account)
	require.NotNil(t, account.PrivateKey)
	require.NotNil(t, account.PublicKey)
}

func TestAccountFromMnemonicBip44_12Words(t *testing.T) {
	mnemonic, err := common.NewMnemonic(common.MnemonicLength12)
	require.NoError(t, err)
	require.NotEmpty(t, mnemonic)

	account, err := common.DeriveAccountFromMnemonicBip44(mnemonic)
	require.NoError(t, err)
	require.NotNil(t, account)
	require.NotNil(t, account.PrivateKey)
	require.NotNil(t, account.PublicKey)
}

func TestAccountFromMnemonicBip39_24Words(t *testing.T) {
	mnemonic, err := common.NewMnemonic(common.MnemonicLength24)
	require.NoError(t, err)
	require.NotEmpty(t, mnemonic)

	account, err := common.DeriveAccountFromMnemonicBip39(mnemonic)
	require.NoError(t, err)
	require.NotNil(t, account)
	require.NotNil(t, account.PrivateKey)
	require.NotNil(t, account.PublicKey)
}

func TestAccountFromMnemonicBip44_24Words(t *testing.T) {
	mnemonic, err := common.NewMnemonic(common.MnemonicLength24)
	require.NoError(t, err)
	require.NotEmpty(t, mnemonic)

	account, err := common.DeriveAccountFromMnemonicBip44(mnemonic)
	require.NoError(t, err)
	require.NotNil(t, account)
	require.NotNil(t, account.PrivateKey)
	require.NotNil(t, account.PublicKey)
}

func TestAccountBase58(t *testing.T) {
	acc := types.NewAccount()

	base58 := common.AccountToBase58(acc)
	require.NotEmpty(t, base58)

	account2, err := common.AccountFromBase58(base58)
	require.NoError(t, err)
	require.NotNil(t, account2)
	require.NotNil(t, account2.PrivateKey)
	require.NotNil(t, account2.PublicKey)
	require.Equal(t, acc, account2)
}

func TestValidateSolanaWalletAddr(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid public key",
			args: args{
				addr: types.NewAccount().PublicKey.ToBase58(),
			},
			wantErr: false,
		},
		{
			name: "invalid public key",
			args: args{
				addr: "invalid",
			},
			wantErr: true,
		},

		{
			name: "invalid public key: too long",
			args: args{
				addr: fmt.Sprintf("%s%s", types.NewAccount().PublicKey.ToBase58(), "q"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := common.ValidateSolanaWalletAddr(tt.args.addr); (err != nil) != tt.wantErr {
				t.Errorf("ValidateSolanaWalletAddr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
