package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

var token = flag.String("token", "", "TelegramBot token key")

func main() {
	flag.Parse()

	pref := tele.Settings{
		Token:  *token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle(tele.OnText, func(c tele.Context) error {
		go kwr(c)

		return nil
	})

	b.Start()
}

func kwr(c tele.Context) error {
	if !c.Message().Private() { // Only reply to private chat
		return nil
	}
	msg := c.Text()
	re := regexp.MustCompile(`^.r([\d]*)d(\d*)[\s]?(.*)`)

	if !re.MatchString(msg) {
		return nil
	}

	ss := re.FindAllSubmatch([]byte(msg), -1)[0]

	var sc int = 3
	if len(ss[1]) == 0 && len(ss[2]) == 0 {
		sc = 0
	} else if len(ss[2]) == 0 {
		sc = 1
	} else if len(ss[1]) == 0 {
		sc = 2
	}

	var res string
	var randnl []int
	switch sc { // Generate random number array
	case 0:
		randnl = get_rand_numa(100)
	case 1:
		d := get_int(string(ss[1]))
		randnl = get_rand_numa(100, d)
	case 2:
		d := get_int(string(ss[2]))
		randnl = get_rand_numa(d)
	case 3:
		d1 := get_int(string(ss[1]))
		d2 := get_int(string(ss[2]))
		randnl = get_rand_numa(d2, d1)
	}

	if len(randnl) > 1 {
		res = fmt.Sprintf("您的结果为: \n%v\n", get_strArray(randnl))
	} else if len(randnl) == 1 {
		res = fmt.Sprintf("您的结果为: %v\n", get_strArray(randnl))
	} else {
		res = "输入有误！"
	}

	if len(ss[3]) != 0 {
		res = fmt.Sprintf("针对 <code>%s</code> ，%s\n", string(ss[3]), res)
	}

	return c.Send(res, tele.ModeHTML)
}

func get_rand_numa(d int, m ...int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var n []int
	if d < 1 { // Wrong Dice, e.g. rd0
		return n
	}
	if len(m) > 0 {
		if m[0] < 1 { // Wrong Dice, e.g. r0d
			return n
		}
		for i := 0; i < m[0]; i++ {
			n = append(n, r.Intn(d)+1)
		}
	} else {
		n = append(n, r.Intn(d)+1)
	}

	return n
}

func get_int(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func get_strArray(ar []int) string {
	return strings.Trim(fmt.Sprint(ar), "[]")
}
