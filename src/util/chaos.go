package util

import (
	"bytes"
	"carrot-backyard/param"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
	"unicode"
)

type hitokotoResponse struct {
	Hitokoto string `json:"hitokoto"`
	From     string `json:"from"`
}

var (
	defaultSentence = hitokotoResponse{
		Hitokoto: "嗨呀",
		From:     "Carrot卡洛塔",
	}
)

func getHitokotoSentence() hitokotoResponse {
	url := "https://v1.hitokoto.cn"
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(nil))
	if err != nil {
		return defaultSentence
	}

	q := req.URL.Query()
	q.Add("c", "a")
	q.Add("c", "b")
	q.Add("c", "c")
	q.Add("encode", "json")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return defaultSentence
	}

	body, _ := ioutil.ReadAll(resp.Body)
	hitokotoResp := hitokotoResponse{}
	if err = json.Unmarshal(body, &hitokotoResp); err != nil {
		return defaultSentence
	}
	return hitokotoResp
}

func GetHitokotoWarpedMessage(message string) string {
	hitokotoResp := getHitokotoSentence()
	return fmt.Sprintf("「%s」\n%s\nfrom.%s", hitokotoResp.Hitokoto, message, hitokotoResp.From)
}

func getMessageChaosVersion(message string) string {
	messageRune := []rune(message)
	messageLen := len(messageRune)
	for i := 0; i <= messageLen-1 && i <= messageLen/9; i++ {
		a, err := rand.Int(rand.Reader, big.NewInt(int64(messageLen-1)))
		wordi, wordii := messageRune[a.Int64()], messageRune[a.Int64()+1]
		if err != nil ||
			unicode.IsNumber(wordi) || unicode.IsNumber(wordii) ||
			unicode.IsLetter(wordi) || unicode.IsLetter(wordii) {
			continue
		}
		messageRune[a.Int64()], messageRune[a.Int64()+1] = wordii, wordi
	}
	return string(messageRune)
}

func getMessageLinkMixedVersion(message string) string {
	messageRune := []rune(message)
	messageLen := len(messageRune)
	messageNew := ""
	for i := 0; i < messageLen; i++ {
		wordi := messageRune[i]
		if wordi == '.' || wordi == ':' {
			messageNew = fmt.Sprintf("%s%s", messageNew, GetRandomEmojiCQString())
			continue
		}
		messageNew = fmt.Sprintf("%s%s", messageNew, string(messageRune[i]))
	}
	return messageNew
}

// SendSameMessageToManyFriends : 批量发送同一条消息，混淆汉字顺序和添加无关内容后不均匀延迟发送
func SendSameMessageToManyFriends(message string, people []param.PersonWithQQ) []param.PersonWithQQ {
	var failed []param.PersonWithQQ
	timer := time.NewTimer(time.Duration(20))
	for _, person := range people {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(20)))
		if err != nil {
			num = big.NewInt(20)
		}
		timer.Reset(time.Second * time.Duration(num.Int64()))
		select {
		case <-timer.C:
			message = getMessageLinkMixedVersion(message)
			message = getMessageChaosVersion(message)
			message = GetHitokotoWarpedMessage(message)

			status := QQSendAndFindWhetherSuccess(person.QQ, message)
			if status == false {
				failed = append(failed, person)
			}
		}
	}
	return failed
}

var (
	emojiInvalid       map[int64]bool
	emojiMessageString []string
)

// packageMessage 在消息后面增加一个随机表情
func packageMessage(message string) string {
	return fmt.Sprintf("%s%s", message, GetRandomEmojiCQString())
}

func GetRandomEmojiCQString() string {
	emoji, err := rand.Int(rand.Reader, big.NewInt(int64(len(emojiMessageString))))
	if err != nil {
		return "👻"
	}
	return emojiMessageString[emoji.Int64()]
}

func buildValidEmoji() {
	emojiInvalid = make(map[int64]bool)
	invalid := []int64{17,
		40, 44, 45, 47, 48,
		51, 52, 58,
		62, 65, 68,
		70, 71, 72, 73,
		80, 82, 83, 84, 87, 88,
		90, 91, 92, 93, 94, 95,
		139,
		141, 142, 143, 149,
		150, 152, 153, 154, 155, 156, 157, 159,
		160, 161, 162, 163, 164, 165, 166, 167,
		170,
		251, 252, 253, 254, 255,
	}
	for _, e := range invalid {
		emojiInvalid[e] = true
	}
	for i := 0; i <= 340; i++ {
		_, exist := emojiInvalid[int64(i)]
		if exist {
			continue
		}
		emojiMessageString = append(emojiMessageString, fmt.Sprintf("[CQ:face,id=%d]", i))
	}

	message := "🙈🙉🙊💘💔💯💤😁😂😃😄👿😉😊😌😍😏😒😓😔😖😘😚😜😝😞😠😡😢😣😥😨😪😭😰😱😲😳😷" +
		"🙃😋😗😛🤑🤓😎🤗🙄🤔😩😤🤐🤒😴😀😆😅😇🙂😙😟😕🙁😫😶😐😑😯😦😧😮😵😬🤕😈👻\U0001F97A\U0001F974" +
		"🤣\U0001F970🤩🤤🤫🤪🧐🤬🤧🤭🤠🤯🤥\U0001F973🤨🤢🤡🤮\U0001F975\U0001F976💩💀👽👾👺👹🤖😺" +
		"😸😹😻😼😽🙀😿😾"
	messageRune := []rune(message)
	for i := range messageRune {
		emojiMessageString = append(emojiMessageString, string(messageRune[i]))
	}
}

func init() {
	buildValidEmoji()
}
