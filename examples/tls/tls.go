package tls

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	cert = "examples/tls/cert.pem"
	key  = "examples/tls/key.pem"
)

type Req struct {
	data []byte
	done chan bool
}

type MyListener chan *MyConn

var ln = MyListener(make(chan *MyConn))

func (ln MyListener) Accept() (c net.Conn, err error) {
	return <-ln, nil
}

func (ln MyListener) Close() error {
	return nil
}

func (ln MyListener) Addr() net.Addr {
	return &net.TCPAddr{net.IP{127, 0, 0, 1}, 49706, ""}
}

type MyConn struct {
	data []byte
	done chan bool
}

func (c *MyConn) Read(b []byte) (n int, err error) {
	if len(c.data) == 0 {
		return 0, io.EOF
	}
	n = copy(b, c.data)
	c.data = c.data[n:]
	return
}

func (c *MyConn) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (c *MyConn) Close() error {
	close(c.done)
	return nil
}

func (c *MyConn) LocalAddr() net.Addr {
	return &net.TCPAddr{net.IP{127, 0, 0, 1}, 49706, ""}
}

func (c *MyConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{net.IP{127, 0, 0, 1}, 49706, ""}
}

func (c *MyConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *MyConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *MyConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func init() {
	go func() {
		c, err := ioutil.ReadFile(cert)
		if err != nil {
			panic(err)
		}
		k, err := ioutil.ReadFile(key)
		if err != nil {
			panic(err)
		}
		cert, err := tls.X509KeyPair(c, k)
		if err != nil {
			panic(err)
		}
		tlsConfig := &tls.Config{
			NextProtos:   []string{"http/1.1"},
			Certificates: []tls.Certificate{cert},
		}
		tlsListener := tls.NewListener(ln, tlsConfig)
		http.HandleFunc("/", handler)
		if err := http.Serve(tlsListener, nil); err != nil {
			panic(err)
		}
	}()
}

var reply = []byte("hello")

func handler(w http.ResponseWriter, req *http.Request) {
	w.Write(reply)
}

func Fuzz(data []byte) int {
	done := make(chan bool)
	ln <- &MyConn{data, done}
	<-done
	return 0
}

/*
package main

import (
	"crypto/tls"
	"encoding/hex"
	//"fmt"
	//"io/ioutil"
	"net"
	"net/http"
	"time"
)

func main() {
	addr := "localhost:49706"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	go func() {
		c, _ := hex.DecodeString(cert)
		k, _ := hex.DecodeString(key)
		cert, err := tls.X509KeyPair(c, k)
		if err != nil {
			panic(err)
		}
		tlsConfig := &tls.Config{
			NextProtos:   []string{"http/1.1"},
			Certificates: []tls.Certificate{cert},
		}
		tlsListener := tls.NewListener(ln, tlsConfig)
		http.HandleFunc("/", handler)
		panic(http.Serve(tlsListener, nil))
	}()
	for i := 0; i < 50; i++ {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	data0 := []byte{0x16, 0x03, 0x1a, 0x00, 0x00}
	c.Write(data0)
	c.Close()
	//c.(*net.TCPConn).CloseWrite()
	//data, err := ioutil.ReadAll(c)
	//fmt.Printf("read returned %q (%v)\n", data, err)
	}
	time.Sleep(time.Second)
	panic("aaa")
}

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("hello"))
}

const cert = "2d2d2d2d2d424547494e2043455254494649434154452d2d2d2d2d0a4d494943" +
	"2b6a434341654b6741774942416749524149696670543942422f416f6a454a58" +
	"4c705558527638774451594a4b6f5a496876634e4151454c425141770a456a45" +
	"514d4134474131554543684d4851574e745a534244627a4165467730784e5441" +
	"304d546b784e544d334e544a61467730784e6a41304d5467784e544d330a4e54" +
	"4a614d4249784544414f42674e5642416f544230466a62575567513238776767" +
	"45694d4130474353714753496233445145424151554141344942447741770a67" +
	"67454b416f4942415144367643635674466c6f72436f455832334d56354f5736" +
	"4b41625a5a4c34467732416f2b6252586d66315363502f445a43704f3464750a" +
	"4e4e3547446a757930484f6a3659495149505668342b4b6b6f4154784e507351" +
	"6e4135744c774b516663616e4a6f4c626f4a6f7a47347a543142427166666365" +
	"0a71346f456769347a3957776b473966747a78424d4a4341497a6c626d2f2f38" +
	"35546f36484f666154686f375577453632334a732b4e3056515449694f593367" +
	"2f0a6348335a7631544338474b577a526930383578794957756c48666545654b" +
	"744f556f76384a7071423365577951656b474f657847532b414141526a754c38" +
	"34510a742f46477a72507970356f4e314d717072377755735472343071307034" +
	"564c653134414142627342656c61432f643132716945736939634a614f796a59" +
	"4148550a46662b2f2f4a302f4d4f374779794d314f444979734b637831614f33" +
	"6e622b6841674d424141476a537a424a4d41344741315564447745422f775145" +
	"417749460a6f44415442674e56485355454444414b4267677242674546425163" +
	"444154414d42674e5648524d4241663845416a41414d42514741315564455151" +
	"4e4d4175430a43577876593246736147397a6444414e42676b71686b69473977" +
	"30424151734641414f4341514541746b4b6d345548334e434749756a75714277" +
	"59416a3052690a31435539397044647749686e697265454b536b344568767451" +
	"75504279696b5951526635434b2f68774c624b6272393836616a6633664a4e43" +
	"7a53656b53614e0a4c542f4e48394743765676744a464a79642f3241775a694f" +
	"5830594843423241362b4b5243794b57524c2b7437774477336175537577634c" +
	"45315a5a6e6d4e690a717673766e6f5072526e3633374c5a3865686c6c437830" +
	"64397068486653302f52435451506b34434451553335496b4235596642766d70" +
	"6430367675734a466d0a65567258322f4b6a3073464376642b70713374345357" +
	"702f4b6d52304769676c514951507542384538663574654341453838334f6b67" +
	"6b2b304b64384339392f0a6e455452764f366242722b4f6277636b4b37424834" +
	"736f41766862323755494452666654433334573470502f635a523848414f6749" +
	"5a7a34334e6b4838413d3d0a2d2d2d2d2d454e44204345525449464943415445" +
	"2d2d2d2d2d0a"

const key = "2d2d2d2d2d424547494e205253412050524956415445204b45592d2d2d2d2d0a" +
	"4d4949457041494241414b43415145412b72776e4662525a614b777142463974" +
	"7a4665546c756967473257532b42634e674b506d3056356e39556e442f773251" +
	"0a71547548626a5465526734377374427a6f2b6d434543443159655069704b41" +
	"4538545437454a774f625338436b48334770796143323643614d78754d303951" +
	"510a616e33334871754b424949754d2f56734a4276583763385154435167434d" +
	"355735762f2f4f55364f687a6e326b34614f314d424f74747962506a64465545" +
	"79490a6a6d4e345033423932623955777642696c73305974504f636369467270" +
	"52333368486972546c4b4c2f4361616764336c736b4870426a6e73526b766741" +
	"4145590a37692f4f454c66785273367a387165614464544b71612b38464c4536" +
	"2b4e4b744b6546533374654141415737415870576776336464716f684c497658" +
	"43576a730a6f3241423142582f762f7964507a44757873736a4e5467794d7243" +
	"6e4d64576a7435322f6f51494441514142416f494241456f75374f6c78434a72" +
	"3968526a790a47777841525078374b784f774138472f494647564c4b39355851" +
	"796e6e494f54776438774b3675686e4c6e68634378426e465538757979476339" +
	"4e596c792f5a0a344678314c6d73466f51635178354e306778666c5077593443" +
	"43646932564737426e686265496673664e4e45714f4c616b2f70432b716e7a66" +
	"34446e6c3072440a73645a366d3071625158516b354231567a47434a33554168" +
	"475256597377364857476e6c5551774d4b476f684152584c5955557436367476" +
	"79346c556c64374e0a5655794b3651624d57635135664769356361454f6a5139" +
	"4e46456c744a2b524437582b697a6d53566f4878484752736a73717268744d7a" +
	"5371333943715264380a58346250476f6658464465373378694e774333514747" +
	"695573663752394e4b61476f745576457548796d6b5644704636786d44454b6e" +
	"4c31305a73524a3978390a3257597a50746b43675945412f363042547137526f" +
	"7343762f723978423848572b5142615332765942386f71534538704c43717542" +
	"6955702f364365506a50650a634e644c457a5a722b48597a71426b52426d4f33" +
	"7835314c4277703338614e336f3350793062365857742f476773776647336372" +
	"505a6d4f4f51336f61754c2f0a35554673635166647a582b515353432f363848" +
	"6341454f6248576b6550382f3268564e4866785639356235555477594c394b67" +
	"4c45717343675945412b77324c0a4d5a3554574b39304e49352f524565737473" +
	"4e663653325a744d3251704a6b646c74764954503943354c504b4448702b4b38" +
	"2b395571484163444f744e4e6a430a72306e30336f3933756942617a714d7932" +
	"4247675236323173747065374c37715063516f4b626a72704a68495133636349" +
	"36465062774c346355385639754f680a672f327a486b4e305668324638434b32" +
	"7a7a6b362b2b466f56445339656a3057554571556c754d4367594541334e4630" +
	"345a6f48494d4f2b7651786d2b4d596e0a66714d5a575335705245454d783672" +
	"6d366d68714b4a434d6432556e68703252726d2b6a505a4b784b63516331542f" +
	"672f6c3239616a2b6c667731426a6f63610a57796458506d4f586f547248336b" +
	"7568536a31674d5447674c684b652b3048577452414f4d6f6b53766474417149" +
	"674b65666e536f722f42426d4f315a6e4f630a664958797141584e3246444c79" +
	"2f787938766a337030554367594541355a454537324a767048454d4f654c7a7a" +
	"5751644d794b45325a79784b50757767464c6a0a4538663136544b68344b2f6d" +
	"326e4a49585a656737366170616642584f6a5063457033324a47336364583651" +
	"69745141386b4e7235522b6250756b67564378660a3167744244715869464b69" +
	"4c712b57472f61334d44523853503871707378474436455a64504264436b6c78" +
	"38315a466f7953543049732b44727a78713578526c0a4378616e755445436759" +
	"4145737577484d524332303337484d4f6e314641314b49626a746139616b6867" +
	"455a75624e5830464271553452764e482b32376c354b0a7a5537796563373547" +
	"65626664355a6f4e482b6341782f6c72562b56422f71693971424b3356584d74" +
	"596a4962794e51596d434d6b4c51476368564d575741750a6e6e755a4f465471" +
	"643358766a6a5a364358465742593970764939476c5163455a4343665a767752" +
	"69696f76306f436a674d53646e513d3d0a2d2d2d2d2d454e4420525341205052" +
	"4956415445204b45592d2d2d2d2d0a"
*/