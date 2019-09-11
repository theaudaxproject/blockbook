// +build unittest

package audax

import (
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"
	"encoding/hex"
	"math/big"
	"os"
	"reflect"
	"testing"

	"github.com/martinboehm/btcutil/chaincfg"
)

func TestMain(m *testing.M) {
	c := m.Run()
	chaincfg.ResetParams()
	os.Exit(c)
}

func Test_GetAddrDescFromAddress_Mainnet(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "P2PKH1",
			args:    args{address: "AGhtPfzWXejd5SVnzUnSdzuXcEQb1qpKqg"},
			want:    "76a9140a371554b9f5958e129f8f276987526e4e1b627f88ac",
			wantErr: false,
		},
		{
			name:    "P2PKH2",
			args:    args{address: "QPWHQu9AAcRnvGeuj6guDkyQ5itkCpJtDE"},
			want:    "76a9141fd45366a2ddcd66cc48ea7f535b449b2a80a76188ac",
			wantErr: false,
		},
		{
			name:    "P2SH1",
			args:    args{address: "bc1q0v3tadxj6pm3ym9j06v9rfyw0jeh5f8squ3nvt"},
			want:    "00147b22beb4d2d077126cb27e9851a48e7cb37a24f0",
			wantErr: false,
		},
		{
			name:    "P2SH2",
			args:    args{address: "bc1qumpyvyxz25kfjjrvyxn3zlyc2wfc0m3l3gm5pg99c4mxylemfqhsdf5q0k"},
			want:    "0020e6c24610c2552c99486c21a7117c98539387ee3f8a3740a0a5c576627f3b482f",
			wantErr: false,
		},
	}
	parser := NewAudaxParser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.GetAddrDescFromAddress(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddrDescFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("GetAddrDescFromAddress() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_GetAddressesFromAddrDesc(t *testing.T) {
	type args struct {
		script string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		want2   bool
		wantErr bool
	}{
		{
			name:    "P2PKH",
			args:    args{script: "76a9140a371554b9f5958e129f8f276987526e4e1b627f88ac"},
			want:    []string{"AGhtPfzWXejd5SVnzUnSdzuXcEQb1qpKqg"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "P2SH",
			args:    args{script: "76a9141fd45366a2ddcd66cc48ea7f535b449b2a80a76188ac"},
			want:    []string{"QPWHQu9AAcRnvGeuj6guDkyQ5itkCpJtDE"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "P2WPKH",
			args:    args{script: "00147b22beb4d2d077126cb27e9851a48e7cb37a24f0"},
			want:    []string{"bc1q0v3tadxj6pm3ym9j06v9rfyw0jeh5f8squ3nvt"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "P2WSH",
			args:    args{script: "0020e6c24610c2552c99486c21a7117c98539387ee3f8a3740a0a5c576627f3b482f"},
			want:    []string{"bc1qumpyvyxz25kfjjrvyxn3zlyc2wfc0m3l3gm5pg99c4mxylemfqhsdf5q0k"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "OP_RETURN ascii",
			args:    args{script: "6a0461686f6a"},
			want:    []string{"OP_RETURN (ahoj)"},
			want2:   false,
			wantErr: false,
		},
		{
			name:    "OP_RETURN hex",
			args:    args{script: "6a072020f1686f6a20"},
			want:    []string{"OP_RETURN 2020f1686f6a20"},
			want2:   false,
			wantErr: false,
		},
	}

	parser := NewAudaxParser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.script)
			got, got2, err := parser.GetAddressesFromAddrDesc(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("outputScriptToAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

var (
	testTx1 bchain.Tx

	testTxPacked1 = "0a20f05ba72a05c4900ff2a00a0403697750201e41267aeea8a589a7dc7bcc57076e12d30101000000010c396f3768565c707addf85ecf47e04cefb2721d95afd977e13f25904de8336a0100000049483045022100fe1b79f38ca4b9dc2fbc50eac6c9bf050ae5a3ee37da05b950918230eb0a8c7c0220477173b60ec00a8b4b28d5db9fc5d8259eea35f91bbdd60089c5ec1be7b05f3b01ffffffff03000000000000000000dd400def140000002321025d145b77df04c40ceb88ea36828755f8275dc5fecf19d3ecccce2d8198c3407cac00d2496b000000001976a914afe70b2e1bf4199298ed8281767bae22970b415088ac0000000018e6b092e605200028f48d1b32770a0012206a33e84d90253fe177d9af951d72b2ef4ce047cf5ef8dd7a705c5668376f390c18012249483045022100fe1b79f38ca4b9dc2fbc50eac6c9bf050ae5a3ee37da05b950918230eb0a8c7c0220477173b60ec00a8b4b28d5db9fc5d8259eea35f91bbdd60089c5ec1be7b05f3b0128ffffffff0f3a04100022003a520a0514ef0d40dd10011a2321025d145b77df04c40ceb88ea36828755f8275dc5fecf19d3ecccce2d8198c3407cac22223764396a4e79716835694b555a566e516a6d7a6d4541513878444b777a4536536e673a470a046b49d20010021a1976a914afe70b2e1bf4199298ed8281767bae22970b415088ac22223769536a6e57436f41556d4a347656584d6b61457561736e5658537a5a7166504c324001"
)

func init() {
	testTx1 = bchain.Tx{
		Hex:       "010000000101c337471f82f50dc527ca44d40c16aab9e8a1c351c67fa4f84ca15970b2f799020000006b483045022100b447e7fc8ea95caf221171a2b9f1b808450b50ba354f71dedf1915a2fec816fa0220154d531481ccb4c5c52fe70a58ef0ce15fa6879587cee0de7f51fdf85f2903ac0121022c50dd67570ad857c78851706a58a50beb80433fe13ef596c843f2cba1df2badffffffff03000000000000000000aeff5d5f000000002321022c50dd67570ad857c78851706a58a50beb80433fe13ef596c843f2cba1df2badac0008af2f000000001976a9147b56c6d93e9bb83016fa99fc98e22e6293e2d86e88ac00000000",
		Blocktime: 1567468329,
		Txid:      "cfebd5729cb7c992fb7b7cb4c1a90df0a150235d83a94d32c024eb8b28f3924f",
		LockTime:  0,
		Time:      1567468329,
		Version:   1,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "483045022100b447e7fc8ea95caf221171a2b9f1b808450b50ba354f71dedf1915a2fec816fa0220154d531481ccb4c5c52fe70a58ef0ce15fa6879587cee0de7f51fdf85f2903ac0121022c50dd67570ad857c78851706a58a50beb80433fe13ef596c843f2cba1df2bad",
				},
				Txid:     "99f7b27059a14cf8a47fc651c3a1e8b9aa160cd444ca27c50df5821f4737c301",
				Vout:     2,
				Sequence: 4294967295,
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(0),
				N:        0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "",
					Addresses: []string{
						"",
					},
				},
			},
			{
				ValueSat: *big.NewInt(1599995822),
				N:        1,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "21022c50dd67570ad857c78851706a58a50beb80433fe13ef596c843f2cba1df2badac",
					Addresses: []string{
						"AT22iSw5YfCeikVQ7RUFGqSPBuhJosCXPr",
					},
				},
			},
			{
				ValueSat: *big.NewInt(800000000),
				N:        2,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "76a9147b56c6d93e9bb83016fa99fc98e22e6293e2d86e88ac",
					Addresses: []string{
						"AT22iSw5YfCeikVQ7RUFGqSPBuhJosCXPr",
					},
				},
			},
		},
	}
}

func Test_PackTx(t *testing.T) {
	type args struct {
		tx        bchain.Tx
		height    uint32
		blockTime int64
		parser    *AudaxParser
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "audax-1",
			args: args{
				tx:        testTx1,
				height:    105509,
				blockTime: 1567468329,
				parser:    NewAudaxParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    testTxPacked1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.parser.PackTx(&tt.args.tx, tt.args.height, tt.args.blockTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("packTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("packTx() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_UnpackTx(t *testing.T) {
	type args struct {
		packedTx string
		parser   *AudaxParser
	}
	tests := []struct {
		name    string
		args    args
		want    *bchain.Tx
		want1   uint32
		wantErr bool
	}{
		{
			name: "audax-1",
			args: args{
				packedTx: testTxPacked1,
				parser:   NewAudaxParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    &testTx1,
			want1:   105509,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.packedTx)
			got, got1, err := tt.args.parser.UnpackTx(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpackTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unpackTx() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("unpackTx() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
