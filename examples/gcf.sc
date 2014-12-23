/* Calculates the greatest common factor. */
var a int
var b int
var f int

a = 45
b = 60
while a != b {
    if a < b {
        b = b - a
    }
    if b < a {
        a = a - b
    }
}
f = a
