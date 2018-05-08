package go_retry_test

import (
	"testing"
	"time"
	. "go-retry"
	"reflect"
	"fmt"
)

type Fail struct{}
type Success struct{
	Data string
}

func (e Fail) Error() string{
	return ""
}

func (e Success) Error() string{
	return ""
}

func Test_DefaultRetryStrategy(t *testing.T){
	retry := NewRetry(&DefaultRetryStrategy{
		TotalTimeout: 5 * time.Second,
		RetryCount: 10,
		RetryErrors: []error{Fail{}},
	})
	test_func := NewTest_function(8)
	err := retry.Do(func() error{
		return test_func()
	})
	if reflect.TypeOf(err) != reflect.TypeOf(Success{}){
		t.Errorf("error=%v\n", reflect.TypeOf(err))
	}
}

func Test_DefaultRetryStrategy_CountOver(t *testing.T){
	retry := NewRetry(&DefaultRetryStrategy{
		TotalTimeout: 5 * time.Second,
		RetryCount: 5,
		RetryErrors: []error{Fail{}},
	})
	test_func := NewTest_function(7)
	err := retry.Do(func() error{
		return test_func()
	})
	if reflect.TypeOf(err) != reflect.TypeOf(CountOverError{}){
		t.Errorf("error=%v\n", reflect.TypeOf(err))
	}
}


func Test_SimpleRetryStrategy(t *testing.T){
	retry := NewRetry(&SimpleRetryStrategy{
		TotalTimeout: 20 * time.Second,
		RetryCount: 10,
		RetryErrors: []error{Fail{}},
		WaitTime: 1 * time.Second,
	})
	test_func := NewTest_function(3)
	err := retry.Do(test_func)
	if reflect.TypeOf(err) != reflect.TypeOf(Success{}){
		t.Errorf("error=%v\n", reflect.TypeOf(err))
	}
}

func Test_SimpleRetryStrategy_CountOver(t *testing.T){
	retry := NewRetry(&SimpleRetryStrategy{
		TotalTimeout: 5 * time.Second,
		RetryCount: 2,
		RetryErrors: []error{Fail{}},
		WaitTime: 1 * time.Second,
	})
	test_func := NewTest_function(3)
	err := retry.Do(func() error{
		return test_func()
	})
	if reflect.TypeOf(err) != reflect.TypeOf(CountOverError{}){
		t.Errorf("error=%v\n", reflect.TypeOf(err))
	}
}

func Test_SimpleRetryStrategy_TotalTimeout(t *testing.T){
	retry := NewRetry(&SimpleRetryStrategy{
		TotalTimeout: 3 * time.Second,
		RetryCount: 5,
		RetryErrors: []error{Fail{}},
		WaitTime: 1 * time.Second,
	})
	test_func := NewTest_function(5)
	err := retry.Do(func() error{
		return test_func()
	})
	if reflect.TypeOf(err) != reflect.TypeOf(TotalTimeoutError{}){
		t.Errorf("error=%v\n", reflect.TypeOf(err))
	}
}




func Test_DefaultRetryStrategy_TotalTimeout(t *testing.T){
	retry := NewRetry(&DefaultRetryStrategy{
		TotalTimeout: 10 * time.Second,
		RetryCount: 5,
		RetryErrors: []error{Fail{}},
	})
	test_func := NewTest_function(3)
	err := retry.Do(func() error{
		time.Sleep(20 * time.Second)
		return test_func()
	})
	if reflect.TypeOf(err) != reflect.TypeOf(TotalTimeoutError{}){
		t.Errorf("error=%v\n", reflect.TypeOf(err))
	}
}

func Test_BackOffRetryStrategy(t *testing.T){
	retry := NewRetry(&BackOffRetryStrategy{
		TotalTimeout: 20 * time.Second,
		RetryCount: 10,
		RetryErrors: []error{Fail{}},
		InitialInterval: 1 * time.Second,
		Multiplier: 2,
		MaxInterval: 8 * time.Second,
	})
	test_func := NewTest_function(3)
	err := retry.Do(test_func)
	if reflect.TypeOf(err) != reflect.TypeOf(Success{}){
		t.Errorf("error=%v\n", reflect.TypeOf(err))
	}
}


func Test_BackOffRetryStrategy_CountOver(t *testing.T){
	retry := NewRetry(&BackOffRetryStrategy{
		TotalTimeout: 20 * time.Second,
		RetryCount: 3,
		RetryErrors: []error{Fail{}},
		InitialInterval: 1 * time.Second,
		Multiplier: 2,
		MaxInterval: 8 * time.Second,
	})
	test_func := NewTest_function(3)
	err := retry.Do(func() error{
		return test_func()
	})
	if reflect.TypeOf(err) != reflect.TypeOf(CountOverError{}){
		t.Errorf("error=%v\n", reflect.TypeOf(err))
	}
}

func Test_BackOffRetryStrategy_MaxInterval(t *testing.T){
	retry := NewRetry(&BackOffRetryStrategy{
		TotalTimeout: 50 * time.Second,
		RetryCount: 10,
		RetryErrors: []error{Fail{}},
		InitialInterval: 1 * time.Second,
		Multiplier: 2,
		MaxInterval: 4 * time.Second,
	})
	test_func := NewTest_function(5)
	err := retry.Do(func() error{
		return test_func()
	})
	if reflect.TypeOf(err) != reflect.TypeOf(Success{}){
		t.Errorf("error=%v\n", reflect.TypeOf(err))
	}
}

func Test_BackOffRetryStrategy_TotalTimeout(t *testing.T){
	retry := NewRetry(&BackOffRetryStrategy{
		TotalTimeout: 5 * time.Second,
		RetryCount: 10,
		RetryErrors: []error{Fail{}},
		InitialInterval: 1 * time.Second,
		Multiplier: 2,
		MaxInterval: 8 * time.Second,
	})
	test_func := NewTest_function(5)
	err := retry.Do(func() error{
		return test_func()
	})
	if reflect.TypeOf(err) != reflect.TypeOf(TotalTimeoutError{}){
		t.Errorf("error=%v\n", reflect.TypeOf(err))
	}
}

func NewTest_function(max_count int) func() error{
	var count int = 0
	return func() error{
		count++
		fmt.Printf("[%v]%v\n", count, time.Now().Format(time.StampMicro))
		if count > max_count{
			return Success{"success"}
		}else{
			return Fail{}
		}
	}
}

