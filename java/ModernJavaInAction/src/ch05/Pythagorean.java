package ch05;

import java.util.function.IntSupplier;
import java.util.stream.IntStream;
import java.util.stream.Stream;

public class Pythagorean {
	public static void main(String[] args) {
		Stream<int[]> pythagoreanTriples =
			    IntStream.rangeClosed(1, 100).boxed()
			             .flatMap(a ->
			                IntStream.rangeClosed(a, 100)
			                         .filter(b -> Math.sqrt(a*a + b*b) % 1 == 0)
			                         .mapToObj(b ->
			                            new int[]{a, b, (int)Math.sqrt(a * a + b * b)})
			);

		pythagoreanTriples
        .forEach(t ->
        	System.out.println(t[0] + ", " + t[1] + ", " + t[2]));
		
		Stream.iterate(new int[]{0, 1},
	               t -> new int[]{t[1], t[0]+t[1]})
		  .limit(20)
	      .forEach(t -> System.out.println("(" + t[0] + "," + t[1] +")"));
		
		Stream.iterate(new int[]{0, 1},
	               t -> new int[]{t[1],t[0] + t[1]})
	      .limit(10)
	      .map(t -> t[0])
	      .forEach(System.out::println);
		
		IntStream.iterate(0, n -> n < 100, n -> n + 4)
        .forEach(System.out::println);

		IntSupplier fib = new IntSupplier(){
		    private int previous = 0;
		    private int current = 1;
		    public int getAsInt(){
		        int oldPrevious = this.previous;
		        int nextValue = this.previous + this.current;
		        this.previous = this.current;
		        this.current = nextValue;
		        return oldPrevious;
		} };
		IntStream.generate(fib).limit(10).forEach(System.out::println);
		
		
		
	}
}
