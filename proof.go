package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

// 10진수를 16진수로
func IntToHex(n int64) []byte {
	return []byte(strconv.FormatInt(n, 16))
}

// 채굴 난이도
const targetBits = 17

// 오버플로우 방지
const maxNonce = math.MaxInt64

type ProofOfWork struct {
	block  *Block
	target *big.Int // 요구사항의 또다른 말
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}
	return pow
}

// 데이터 준비
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PreBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int // hash의 정수 표현 값
	var hash [32]byte
	nonce := 0 // 카운터

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)

	for nonce < maxNonce {
		data := pow.prepareData(nonce) // 데이터 준비 (생성)
		hash = sha256.Sum256(data)     // SHA-256 해싱
		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])          // 해시값의 큰 정수로의 변환
		if hashInt.Cmp(pow.target) == -1 { // 정수값과 타겟값 비교
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// 작업 증명 검증 기능
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}
