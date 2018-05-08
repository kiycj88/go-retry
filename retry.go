package go_retry

import (
	"reflect"
	"time"
)

type RetryStrategy interface {
	Init()
	Total_Timeout() time.Duration
	NeedRetry(err error) (bool, error)
}

type TotalTimeoutError struct{}

type CountOverError struct{}

func (err TotalTimeoutError) Error() string{
	return ""
}

func (err CountOverError) Error() string{
	return ""
}

type DefaultRetryStrategy struct {
	TotalTimeout 	time.Duration
	RetryCount	int
	RetryErrors	[]error
	attempt		int
}

func (s *DefaultRetryStrategy) Init(){
	if s.TotalTimeout <= 0{
		s.TotalTimeout = 30 * time.Second
	}
	if s.RetryCount <= 0{
		s.RetryCount = 5
	}
}

func (s *DefaultRetryStrategy) Total_Timeout() 	time.Duration{
	return s.TotalTimeout
}

func (s *DefaultRetryStrategy) NeedRetry(err error) (bool, error){
	return func() (bool, error){
		if s.attempt < s.RetryCount{
			errType := reflect.TypeOf(err)
			for _, e := range s.RetryErrors{
				if errType == reflect.TypeOf(e){
					s.attempt++
					return true, err
				}
			}
			return false, err

		}
		return false, CountOverError{}
	}()
}

type SimpleRetryStrategy struct {
	TotalTimeout 	time.Duration
	RetryCount	int
	RetryErrors	[]error
	WaitTime	time.Duration
	attempt		int
}

func (s *SimpleRetryStrategy) Init() {
	if s.TotalTimeout <= 0{
		s.TotalTimeout = 30 * time.Second
	}
	if s.RetryCount <= 0{
		s.RetryCount = 5
	}
	if s.WaitTime <= 0 {
		s.WaitTime = 0
	}
}

func (s *SimpleRetryStrategy) Total_Timeout() 	time.Duration{
	return s.TotalTimeout
}

func (s *SimpleRetryStrategy) NeedRetry(err error) (bool, error){
	return func() (bool, error){
		if s.attempt < s.RetryCount{
			errType := reflect.TypeOf(err)
			for _, e := range s.RetryErrors{
				if errType == reflect.TypeOf(e){
					if s.WaitTime != 0{
						time.Sleep(s.WaitTime)
					}
					s.attempt++
					return true, err
				}
			}
			return false, err

		}
		return false, CountOverError{}
	}()
}

type BackOffRetryStrategy struct {
	TotalTimeout 	time.Duration
	RetryCount	int
	RetryErrors	[]error
	InitialInterval  time.Duration
	MaxInterval	time.Duration
	Multiplier	float32
	attempt		int
	interval	time.Duration
}

func (s *BackOffRetryStrategy) Init(){
	if s.TotalTimeout <= 0{
		s.TotalTimeout = 30 * time.Second
	}
	if s.RetryCount <= 0{
		s.RetryCount = 5
	}
	if s.InitialInterval <= 0 {
		s.InitialInterval = 1 * time.Second
	}
	if s.MaxInterval <= 0 {
		s.MaxInterval = 8 * time.Second
	}
	if s.Multiplier <= 0.0 {
		s.Multiplier = 2.0
	}
	s.interval = s.InitialInterval
}

func (s *BackOffRetryStrategy) Total_Timeout() time.Duration{
	return s.TotalTimeout
}

func (s *BackOffRetryStrategy) NeedRetry(err error) (bool, error){
	return func() (bool, error){
		if s.attempt < s.RetryCount{
			errType := reflect.TypeOf(err)
			for _, e := range s.RetryErrors{
				if errType == reflect.TypeOf(e){
					time.Sleep(s.interval)
					s.interval = time.Duration(float32(s.interval) * s.Multiplier)
					if s.interval > s.MaxInterval{
						s.interval = s.MaxInterval
					}
					s.attempt++
					return true, err
				}
			}
			return false, err

		}
		return false, CountOverError{}
	}()
}

type Retry struct{
	Strategy RetryStrategy
}

func (r *Retry) Do(f func() error) error{
	total_chan := make(chan error)
	stop_chan := make(chan error)
	defer close(stop_chan)

	go func(c chan error){
		defer close(c)
		for{
			timeout := false
			err := f()
			select {
			case <- stop_chan:
				timeout = true
				break
			default:
			}
			if timeout{
				break
			}
			needRetry, result := r.Strategy.NeedRetry(err)
			if !needRetry{
				c <- result
				break
			}
		}
	}(total_chan)

	select{
	case <- time.After(r.Strategy.Total_Timeout()):
		stop_chan <- TotalTimeoutError{}
		return TotalTimeoutError{}
	case e := <- total_chan:
		return e
	}
}

func NewRetry(strategy RetryStrategy) Retry{
	strategy.Init()
	return Retry{strategy}
}
