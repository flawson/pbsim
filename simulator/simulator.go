package simulator

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	MaxWBallValue = 69
	MaxPBallValue = 26
	NumWBalls     = 5
)

type Results struct {
	numAttempts int
	wc1         int
	wc2         int
	wc3         int
	wc4         int
	wc5         int
	wc6         int
	wc7         int
	wc8         int
	wc9         int
}

func RNGenerator(seedval int64) *rand.Rand {
	return rand.New(rand.NewSource(seedval))
}

type Ballset struct {
	wBalls []int
	pBall  int
}

func NewBallset() *Ballset {
	bs := &Ballset{
		wBalls: make([]int, NumWBalls),
	}
	return bs
}

// Checks for matches against passed ballset, returns wb matches followed by
// pb match.
func (bs *Ballset) CheckMatches(bs2 *Ballset) (int, int) {
	pbMatchCnt := 0
	wbMatchCnt := 0

	if bs.pBall == bs2.pBall {
		pbMatchCnt++
	}

	for _, val := range bs.wBalls {
		for _, v2 := range bs2.wBalls {
			if val == v2 {
				wbMatchCnt++
				break
			}
		}
	}

	return wbMatchCnt, pbMatchCnt
}

func Generate(r *rand.Rand) *Ballset {
	bs := NewBallset()
	wBalls := make([]int, MaxWBallValue)
	for i := 0; i < len(wBalls); i++ {
		wBalls[i] = i + 1
	}
	bs.pBall = r.Intn(MaxPBallValue) + 1

	for i := 0; i < NumWBalls; i++ {
		wbIdx := r.Intn(len(wBalls))
		bs.wBalls[i] = wBalls[wbIdx]
		wBalls = append(wBalls[:wbIdx], wBalls[wbIdx+1:]...)
	}
	return bs
}

type Simulator struct {
	Ballset *Ballset
	Results *Results
}

func NewSimulator() *Simulator {
	s := &Simulator{
		Results: &Results{},
	}
	s.Ballset = Generate(RNGenerator(time.Now().UnixMicro()))
	return s
}

func (s *Simulator) Run(numJackpots int) {
	done := make(chan Results)
	for i := 0; i < numJackpots; i++ {
		go func(id int, done chan Results) {
			rng := RNGenerator(time.Now().UnixMicro())
			res := Results{}
			for {
				winSet := Generate(rng)
				res.numAttempts++

				nWballs, pbMatch := s.Ballset.CheckMatches(winSet)

				switch {
				// PB
				case pbMatch == 1 && nWballs == 0:
					res.wc1++
				// PB + 1WB
				case pbMatch == 1 && nWballs == 1:
					res.wc2++
				// PB + 2WB
				case pbMatch == 1 && nWballs == 2:
					res.wc3++
				// 3WB
				case pbMatch == 0 && nWballs == 3:
					res.wc4++
				// PB + 3WB
				case pbMatch == 1 && nWballs == 3:
					res.wc5++
				// 4WB
				case pbMatch == 0 && nWballs == 4:
					res.wc6++
				// PB + 4WB
				case pbMatch == 1 && nWballs == 4:
					res.wc7++
				// 5WB
				case pbMatch == 0 && nWballs == 5:
					res.wc8++
				// PB + 5WB
				case pbMatch == 1 && nWballs == 5:
					res.wc9++
					done <- res
					return
				}
			}
		}(i, done)
	}

	for i := 0; i < numJackpots; i++ {
		if i%2 == 0 {
			fmt.Printf("Pct Complete: %.3f\n", float64(i)/float64(numJackpots)*100)
		}
		v := <-done
		s.Results.numAttempts += v.numAttempts
		s.Results.wc1 += v.wc1
		s.Results.wc2 += v.wc2
		s.Results.wc3 += v.wc3
		s.Results.wc4 += v.wc4
		s.Results.wc5 += v.wc5
		s.Results.wc6 += v.wc6
		s.Results.wc7 += v.wc7
		s.Results.wc8 += v.wc8
		s.Results.wc9 += v.wc9
	}
}

func (s *Simulator) PrintResults() {
	fmt.Printf("Number of attempts: %d\n", s.Results.numAttempts)
	fmt.Printf("WinCond 1: %d, chance of winning: %.3f\n", s.Results.wc1, float64(s.Results.numAttempts)/float64(s.Results.wc1))
	fmt.Printf("WinCond 2: %d, chance of winning: %.3f\n", s.Results.wc2, float64(s.Results.numAttempts)/float64(s.Results.wc2))
	fmt.Printf("WinCond 3: %d, chance of winning: %.3f\n", s.Results.wc3, float64(s.Results.numAttempts)/float64(s.Results.wc3))
	fmt.Printf("WinCond 4: %d, chance of winning: %.3f\n", s.Results.wc4, float64(s.Results.numAttempts)/float64(s.Results.wc4))
	fmt.Printf("WinCond 5: %d, chance of winning: %.3f\n", s.Results.wc5, float64(s.Results.numAttempts)/float64(s.Results.wc5))
	fmt.Printf("WinCond 6: %d, chance of winning: %.3f\n", s.Results.wc6, float64(s.Results.numAttempts)/float64(s.Results.wc6))
	fmt.Printf("WinCond 7: %d, chance of winning: %.3f\n", s.Results.wc7, float64(s.Results.numAttempts)/float64(s.Results.wc7))
	fmt.Printf("WinCond 8: %d, chance of winning: %.3f\n", s.Results.wc8, float64(s.Results.numAttempts)/float64(s.Results.wc8))
	fmt.Printf("WinCond 9: %d, chance of winning: %.3f\n", s.Results.wc9, float64(s.Results.numAttempts)/float64(s.Results.wc9))
}
