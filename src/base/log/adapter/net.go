package adapter

//校验AdapterNet符合Adapter接口
var _ Adapter = &AdapterNet{}
type AdapterNet struct{}

func (an *AdapterNet) Write(b []byte) error {
    return nil
}

func (an *AdapterNet) Close() {

}