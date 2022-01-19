(define real-min -2.5)
(define real-max +1.5)
(define imag-min -2.0)
(define imag-max +2.0)

(define iter-count 40)
(define escape-radius 4.0)

(define width 1024)
(define scale (/ width (- real-max real-min)))
(define height (round->exact (* scale (- imag-max imag-min))))

(define (complex real imag) (make-rectangular real imag))

(define (mandelbrot x y)
  (define (*mandelbrot z0 z rem)
    (if (or (= rem 0) (>= (magnitude z) escape-radius))
        rem
        (*mandelbrot z0 (+ (* z z) z0) (- rem 1))))
  (let ((val (complex (+ (/ x scale) real-min)
                      (+ (/ y scale) imag-min))))
    (round (- 255 (* 255 (/ (- iter-count (*mandelbrot val 0 iter-count)) iter-count))))))

(define (plot file)
  (with-output-to-file file
    (lambda () (begin
      (display "P2") (newline)
      (display width) (newline)
      (display height) (newline)
      (display 255) (newline)
      (do ((y 0 (+ y 1))) ((>= y height))
        (do ((x 0 (+ x 1))) ((>= x width))
          (begin (display (mandelbrot x y)) (newline))))))))

(plot "mandelbrot.pbm")
