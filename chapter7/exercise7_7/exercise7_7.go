// Explain why the help message contains °C when the default value of 20.0 does not.
// Because it uses tempconv.Celsius's String() method to print the help message
// and that method formats the string like so: "%g°C", therefore it passes the
// value (in this case the default) and append "°C" to it.
