package ssh

import (
	"errors"
	"fmt"
	"github.com/DawnBreather/go-commons/file"
	"github.com/DawnBreather/go-commons/logger"
	path2 "github.com/DawnBreather/go-commons/path"
	"github.com/DawnBreather/go-commons/socket"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"net"
)

type SSH struct {
	host        string
	port        int
	user        string
	password    string
	keyPath     string
	keyPassword string

	client *goph.Client
	CMD    *cmd
}

var _logger = logger.New()

func (s *SSH) SetHost(host string) *SSH{
	s.host = host
	return s
}
func (s *SSH) SetPort(port int) *SSH{
	s.port = port
	return s
}
func (s *SSH) SetUsername(user string) *SSH{
	s.user = user
	return s
}
func (s *SSH) SetPassword(password string) *SSH{
	s.password = password
	return s
}
func (s *SSH) SetKeyPath(keyPath string) *SSH{
	s.keyPath = keyPath
	return s
}
func (s *SSH) SetKeyPassword(keyPassword string) *SSH{
	s.keyPassword = keyPassword
	return s
}

func (s *SSH) IsPortOpen() bool {
	st := socket.Socket{}
	return st.
			SetHost(s.host).
			SetPort(s.port).
			SetProtocol(socket.TCP).
			IsPortOpen()
}

func (s *SSH) getSocket() string {
	st := socket.Socket{}
	return st.
		SetHost(s.host).
		SetPort(s.port).
		SetProtocol(socket.TCP).
		GetSocket()
}

func VerifyHost(host string, remote net.Addr, key ssh.PublicKey) error {

	//
	// If you want to connect to new hosts.
	// here your should check new connections public keys
	// if the key not trusted you shuld return an error
	//

	// hostFound: is host in known hosts file.
	// err: error if key not in known hosts file OR host in known hosts file but key changed!
	hostFound, err := goph.CheckKnownHost(host, remote, key, "")

	// Host in known hosts but key mismatch!
	// Maybe because of MAN IN THE MIDDLE ATTACK!
	if hostFound && err != nil {

		return err
	}

	// handshake because public key already exists.
	if hostFound && err == nil {

		return nil
	}

	// Ask user to check if he trust the host public key.
	if askIsHostTrusted(host, key) == false {

		// Make sure to return error on non trusted keys.
		return errors.New("you typed no, aborted!")
	}

	// Add the new host to known hosts file.
	return goph.AddKnownHost(host, remote, key, "")
}

func askIsHostTrusted(host string, key ssh.PublicKey) bool {

	return true

	//reader := bufio.NewReader(os.Stdin)
	//
	//fmt.Printf("Unknown Host: %s \nFingerprint: %s \n", host, ssh.FingerprintSHA256(key))
	//fmt.Print("Would you like to add it? type yes or no: ")
	//
	//a, err := reader.ReadString('\n')
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//return strings.ToLower(strings.TrimSpace(a)) == "yes"
}

func (s *SSH) Connect() *SSH {

	var shouldConnect = false
	if s.client == nil {
		shouldConnect = true
	} else {
		if s.client.Client == nil {
			shouldConnect = true
		}
	}

	if shouldConnect {
		if s.CanConnectWithKey() {
			s.ConnectWithKey()
		} else if s.CanConnectWithPassword() {
			s.ConnectWithPassword()
		}
	}

	return s
}

func (s *SSH) CanConnectWithKey() bool {
	defer s.Close()
	return s.connectWithKey()
}

func (s *SSH) ConnectWithKey() *SSH {
	s.connectWithKey()
	return s
}

func (s *SSH) connectWithKey() bool {
	p := path2.Path{}
	p.SetPath(s.keyPath)
	if p.Exists() {
		if p.IsFile() {
			auth, err := goph.Key(s.keyPath, s.keyPassword)
			if err != nil {
				_logger.Errorf("Unable to read SSH key { %s }: %v", s.keyPath, err)
				return false
			}

			//s.client, err = goph.New(s.user, s.host, auth)
			s.client, err = goph.NewConn(&goph.Config{
				User:     s.user,
				Addr:     s.host,
				Port:     uint(s.port),
				Auth:     auth,
				Callback: VerifyHost,
			})
			if err != nil {
				_logger.Errorf("Unable to establish SSH connection to { %s } with username { %s } and key { %s }: %v", s.getSocket(), s.user, s.keyPath, err)
				return false
			}
		} else {
			return false
		}
	} else {
		return false
	}

	s.CMD = &cmd{
		s: s,
	}

	return true
}

func (s *SSH) CanConnectWithPassword() bool {
	defer s.Close()
	return s.connectWithPassword()
}

func (s *SSH) ConnectWithPassword() *SSH {
	s.connectWithPassword()
	return s
}

func (s *SSH) connectWithPassword() bool {
	var err error
	//s.client, err = goph.New(s.user, s.host, goph.Password(s.password))
	//if err != nil {
	//	_logger.Errorf("Unable to establish SSH connection to { %s } with username { %s } and key { %s }: %v", s.getSocket(), s.user, s.keyPath, err)
	//	return false
	//}
	//
	//s.CMD = &cmd{
	//	s: s,
	//}

	//s.client, err = goph.New(s.user, s.host, auth)
	s.client, err = goph.NewConn(&goph.Config{
		User:     s.user,
		Addr:     s.host,
		Port:     uint(s.port),
		Auth:     goph.Password(s.password),
		Callback: VerifyHost,
	})
	if err != nil {
		_logger.Errorf("Unable to establish SSH connection to { %s } with username { %s } and key { %s }: %v", s.getSocket(), s.user, s.keyPath, err)
		return false
	}

	s.CMD = &cmd{
		s: s,
	}

	return true
}

func (s *SSH) DownloadRemoteFile(fromRemotePath, toLocalPath string) *SSH {
	err := s.client.Download(fromRemotePath, toLocalPath)
	if err != nil {
		_logger.Errorf("Unable to download file from { %s } to { %s } on server { %s@%s }: %v", fromRemotePath, toLocalPath, s.user, s.getSocket(), err)
	}
	return s
}

func (s *SSH) DownloadRemoteFileOverSSHWithSudo(fromRemotePath, toLocalPath string) *SSH {
	f := file.File{}
	p := path2.Path{}
	f.SetPath(toLocalPath)
	p.SetPath(f.GetBaseDir()).MkdirAll(0644)

	content := s.CMD.ExecuteSudoBash(fmt.Sprintf("cat %s", fromRemotePath))

	f.
		SetContent([]byte(content)).
		Save()

	return s
}

func (s *SSH) UploadLocalFile(fromLocalPath, toRemotePath string) *SSH {
	err := s.client.Upload(fromLocalPath, toRemotePath)
	if err != nil {
		_logger.Errorf("Unable to upload file from { %s } to { %s } on server { %s@%s }: %v", fromLocalPath, toRemotePath, s.user, s.getSocket(), err)
	}
	return s
}

func (s *SSH) UploadLocalFileOverSSHWithSudo(fromLocalPath, toRemotePath string) *SSH {
	f := file.File{}
	b64EncodedString := f.SetPath(fromLocalPath).ReadContent().ParseContentToBase64().GetBase64()
	out := s.CMD.ExecuteSudoBash(fmt.Sprintf("echo %s | sudo base64 -d > %s", b64EncodedString, toRemotePath))
	if out != "" {
		_logger.Errorf("Unable to upload local file from { %s } to remote { %s } / { %s }", fromLocalPath, toRemotePath, s.host)
	}
	return s
}

func (s *SSH) Close() *SSH {
	if s.client != nil {
		if s.client.Client != nil {
			err := s.client.Close()
			if err != nil {
				_logger.Errorf("Unable to close SSH connection with { %s }: %v", s.getSocket(), err)
			}
		}
	}
	return s
}

