package net

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestAddMLI(t *testing.T) {

	echoServer := EchoServ{
		TcpAddr: &net.TCPAddr{
			IP:   net.ParseIP("localhost"),
			Port: 8888,
			Zone: "",
		}}
	go func() {
		if err := echoServer.ListenAndAccept(); err != nil {
			t.Fatal(err)

		}
	}()

	t.Run("Test MLI2e", func(t *testing.T) {
		ncc := NewNetCatClient("localhost:8888", Mli2e)
		if err := ncc.OpenConnection(); err != nil {
			t.Fatal(err)
		}
		payload := []byte("hello world")
		withMLI := AddMLI(Mli2e, payload)
		assert.Equal(t, []byte{0, 0x0b}, withMLI[0:2])

		if err := ncc.Write(payload); err != nil {
			t.Fatal(err)
		}
		response, err := ncc.ReadNextPacket()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, payload, response)
	})

	t.Run("Test MLI2i", func(t *testing.T) {
		ncc := NewNetCatClient("localhost:8888", Mli2i)
		if err := ncc.OpenConnection(); err != nil {
			t.Fatal(err)
		}
		payload := []byte("hello world")
		withMLI := AddMLI(Mli2i, payload)
		assert.Equal(t, []byte{0, 0x0d}, withMLI[0:2])

		if err := ncc.Write(payload); err != nil {
			t.Fatal(err)
		}
		response, err := ncc.ReadNextPacket()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, payload, response)
	})

	t.Run("Test MLI4i", func(t *testing.T) {
		ncc := NewNetCatClient("localhost:8888", Mli4i)
		if err := ncc.OpenConnection(); err != nil {
			t.Fatal(err)
		}
		payload := []byte("hello world")
		withMLI := AddMLI(Mli4i, payload)
		assert.Equal(t, []byte{0, 0, 0, 0x0f}, withMLI[0:4])

		if err := ncc.Write(payload); err != nil {
			t.Fatal(err)
		}
		response, err := ncc.ReadNextPacket()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, payload, response)
	})

	t.Run("Test MLI4e", func(t *testing.T) {
		ncc := NewNetCatClient("localhost:8888", Mli4e)
		if err := ncc.OpenConnection(); err != nil {
			t.Fatal(err)
		}

		payload := []byte("hello world")
		withMLI := AddMLI(Mli4e, payload)
		assert.Equal(t, []byte{0, 0, 0, 0x0b}, withMLI[0:4])

		if err := ncc.Write(payload); err != nil {
			t.Fatal(err)
		}
		response, err := ncc.ReadNextPacket()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, payload, response)
	})
}
