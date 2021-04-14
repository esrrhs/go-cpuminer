package main

import "strings"

const (
	INVALID       = iota
	CN_0          // "cn/0"             CryptoNight (original).
	CN_1          // "cn/1"             CryptoNight variant 1 also known as Monero7 and CryptoNightV7.
	CN_2          // "cn/2"             CryptoNight variant 2.
	CN_R          // "cn/r"             CryptoNightR (Monero's variant 4).
	CN_FAST       // "cn/fast"          CryptoNight variant 1 with half iterations.
	CN_HALF       // "cn/half"          CryptoNight variant 2 with half iterations (Masari/Torque).
	CN_XAO        // "cn/xao"           CryptoNight variant 0 (modified Alloy only).
	CN_RTO        // "cn/rto"           CryptoNight variant 1 (modified Arto only).
	CN_RWZ        // "cn/rwz"           CryptoNight variant 2 with 3/4 iterations and reversed shuffle operation (Graft).
	CN_ZLS        // "cn/zls"           CryptoNight variant 2 with 3/4 iterations (Zelerius).
	CN_DOUBLE     // "cn/double"        CryptoNight variant 2 with double iterations (X-CASH).
	CN_LITE_0     // "cn-lite/0"        CryptoNight-Lite variant 0.
	CN_LITE_1     // "cn-lite/1"        CryptoNight-Lite variant 1.
	CN_HEAVY_0    // "cn-heavy/0"       CryptoNight-Heavy (4 MB).
	CN_HEAVY_TUBE // "cn-heavy/tube"    CryptoNight-Heavy (modified TUBE only).
	CN_HEAVY_XHV  // "cn-heavy/xhv"     CryptoNight-Heavy (modified Haven Protocol only).
	CN_PICO_0     // "cn-pico"          CryptoNight-Pico
	CN_PICO_TLO   // "cn-pico/tlo"      CryptoNight-Pico (TLO)
	CN_CCX        // "cn/ccx"           Conceal (CCX)
	RX_0          // "rx/0"             RandomX (reference configuration).
	RX_WOW        // "rx/wow"           RandomWOW (Wownero).
	RX_ARQ        // "rx/arq"           RandomARQ (Arqma).
	RX_SFX        // "rx/sfx"           RandomSFX (Safex Cash).
	RX_KEVA       // "rx/keva"          RandomKEVA (Keva).
	AR2_CHUKWA    // "argon2/chukwa"    Argon2id (Chukwa).
	AR2_CHUKWA_V2 // "argon2/chukwav2"  Argon2id (Chukwa v2).
	AR2_WRKZ      // "argon2/wrkz"      Argon2id (WRKZ)
	ASTROBWT_DERO // "astrobwt"         AstroBWT (Dero)
	KAWPOW_RVN    // "kawpow/rvn"       KawPow (RVN)
	MAX
)

const (
	UNKNOWN = iota
	CN
	CN_LITE
	CN_HEAVY
	CN_PICO
	RANDOM_X
	ARGON2
	ASTROBWT
	KAWPOW
)

type AlgoName struct {
	name      string
	shortName string
	id        int
}

var algorithm_names = []AlgoName{
	{"cryptonight/0", "cn/0", CN_0},
	{"cryptonight", "cn", CN_0},
	{"cryptonight/1", "cn/1", CN_1},
	{"cryptonight-monerov7", "", CN_1},
	{"cryptonight_v7", "", CN_1},
	{"cryptonight/2", "cn/2", CN_2},
	{"cryptonight-monerov8", "", CN_2},
	{"cryptonight_v8", "", CN_2},
	{"cryptonight/r", "cn/r", CN_R},
	{"cryptonight_r", "", CN_R},
	{"cryptonight/fast", "cn/fast", CN_FAST},
	{"cryptonight/msr", "cn/msr", CN_FAST},
	{"cryptonight/half", "cn/half", CN_HALF},
	{"cryptonight/xao", "cn/xao", CN_XAO},
	{"cryptonight_alloy", "", CN_XAO},
	{"cryptonight/rto", "cn/rto", CN_RTO},
	{"cryptonight/rwz", "cn/rwz", CN_RWZ},
	{"cryptonight/zls", "cn/zls", CN_ZLS},
	{"cryptonight/double", "cn/double", CN_DOUBLE},
	{"cryptonight-lite/0", "cn-lite/0", CN_LITE_0},
	{"cryptonight-lite/1", "cn-lite/1", CN_LITE_1},
	{"cryptonight-lite", "cn-lite", CN_LITE_1},
	{"cryptonight-light", "cn-light", CN_LITE_1},
	{"cryptonight_lite", "", CN_LITE_1},
	{"cryptonight-aeonv7", "", CN_LITE_1},
	{"cryptonight_lite_v7", "", CN_LITE_1},
	{"cryptonight-heavy/0", "cn-heavy/0", CN_HEAVY_0},
	{"cryptonight-heavy", "cn-heavy", CN_HEAVY_0},
	{"cryptonight_heavy", "", CN_HEAVY_0},
	{"cryptonight-heavy/xhv", "cn-heavy/xhv", CN_HEAVY_XHV},
	{"cryptonight_haven", "", CN_HEAVY_XHV},
	{"cryptonight-heavy/tube", "cn-heavy/tube", CN_HEAVY_TUBE},
	{"cryptonight-bittube2", "", CN_HEAVY_TUBE},
	{"cryptonight-pico", "cn-pico", CN_PICO_0},
	{"cryptonight-pico/trtl", "cn-pico/trtl", CN_PICO_0},
	{"cryptonight-turtle", "cn-trtl", CN_PICO_0},
	{"cryptonight-ultralite", "cn-ultralite", CN_PICO_0},
	{"cryptonight_turtle", "cn_turtle", CN_PICO_0},
	{"cryptonight-pico/tlo", "cn-pico/tlo", CN_PICO_TLO},
	{"cryptonight/ultra", "cn/ultra", CN_PICO_TLO},
	{"cryptonight-talleo", "cn-talleo", CN_PICO_TLO},
	{"cryptonight_talleo", "cn_talleo", CN_PICO_TLO},
	{"randomx/0", "rx/0", RX_0},
	{"randomx/test", "rx/test", RX_0},
	{"RandomX", "rx", RX_0},
	{"randomx/wow", "rx/wow", RX_WOW},
	{"RandomWOW", "", RX_WOW},
	{"randomx/arq", "rx/arq", RX_ARQ},
	{"RandomARQ", "", RX_ARQ},
	{"randomx/sfx", "rx/sfx", RX_SFX},
	{"RandomSFX", "", RX_SFX},
	{"randomx/keva", "rx/keva", RX_KEVA},
	{"RandomKEVA", "", RX_KEVA},
	{"argon2/chukwa", "", AR2_CHUKWA},
	{"chukwa", "", AR2_CHUKWA},
	{"argon2/chukwav2", "", AR2_CHUKWA_V2},
	{"chukwav2", "", AR2_CHUKWA_V2},
	{"argon2/wrkz", "", AR2_WRKZ},
	{"astrobwt", "", ASTROBWT_DERO},
	{"astrobwt/dero", "", ASTROBWT_DERO},
	{"kawpow", "", KAWPOW_RVN},
	{"kawpow/rvn", "", KAWPOW_RVN},
	{"cryptonight/ccx", "cn/ccx", CN_CCX},
	{"cryptonight/conceal", "cn/conceal", CN_CCX},
}

type Algorithm struct {
	id int
}

func NewAlgorithm(algo string) *Algorithm {
	for _, item := range algorithm_names {
		if strings.EqualFold(algo, item.name) || (item.shortName != "" && strings.EqualFold(algo, item.shortName)) {
			return &Algorithm{item.id}
		}
	}
	return &Algorithm{INVALID}
}

func (a *Algorithm) family() int {

	switch a.id {
	case CN_0:
	case CN_1:
	case CN_2:
	case CN_R:
	case CN_FAST:
	case CN_HALF:
	case CN_XAO:
	case CN_RTO:
	case CN_RWZ:
	case CN_ZLS:
	case CN_DOUBLE:
	case CN_CCX:
		return CN

	case CN_LITE_0:
	case CN_LITE_1:
		return CN_LITE

	case CN_HEAVY_0:
	case CN_HEAVY_TUBE:
	case CN_HEAVY_XHV:
		return CN_HEAVY

	case CN_PICO_0:
	case CN_PICO_TLO:
		return CN_PICO

	case RX_0:
	case RX_WOW:
	case RX_ARQ:
	case RX_SFX:
	case RX_KEVA:
		return RANDOM_X

	case AR2_CHUKWA:
	case AR2_CHUKWA_V2:
	case AR2_WRKZ:
		return ARGON2

	case ASTROBWT_DERO:
		return ASTROBWT
	case KAWPOW_RVN:
		return KAWPOW

	default:
		break
	}

	return UNKNOWN
}

func (a *Algorithm) name() string {
	for _, item := range algorithm_names {
		if item.id == a.id {
			return item.name
		}
	}
	return "invalid"
}

func (a *Algorithm) shortName() string {
	for _, item := range algorithm_names {
		if item.id == a.id {
			return item.shortName
		}
	}
	return "invalid"
}

func (a *Algorithm) supportAlgoName() string {
	switch a.id {
	case CN_0:
		return "cn/0"
	case CN_1:
		return "cn/1"
	case CN_2:
		return "cn/2"
	case CN_R:
		return "cn/r"
	case CN_FAST:
		return "cn/fast"
	case CN_HALF:
		return "cn/half"
	case CN_XAO:
		return "cn/xao"
	case CN_RTO:
		return "cn/rto"
	case CN_RWZ:
		return "cn/rwz"
	case CN_ZLS:
		return "cn/zls"
	case CN_DOUBLE:
		return "cn/double"
	case CN_LITE_0:
		return "cn-lite/0"
	case CN_LITE_1:
		return "cn-lite/1"
	case CN_HEAVY_0:
		return "cn-heavy/0"
	case CN_HEAVY_TUBE:
		return "cn-heavy/tube"
	case CN_HEAVY_XHV:
		return "cn-heavy/xhv"
	case CN_PICO_0:
		return "cn-pico"
	case CN_PICO_TLO:
		return "cn-pico/tlo"
	default:
		break
	}
	return ""
}
