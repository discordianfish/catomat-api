package stepper

import "testing"

func TestNewMotor(t *testing.T) {
	_, err := NewMotor(25, 24, 23, 18)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewMotorError(t *testing.T) {
	_, err := NewMotor(25, 24, 23)
	if err == nil {
		t.Fatal("Expected error due to number of pins")
	}
}
