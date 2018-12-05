package ch04.item23;

public class Rectangle extends Figure2 {
    final double length;
    final double width;
    Rectangle(double length, double width) {
        this.length = length;
        this.width  = width;
    }

    @Override
    double area() {
        return length * width;
    }
}
