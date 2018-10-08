package chapter10;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class ThreadRunnable {
	public static void main(String[] args) {
		Runnable hellos = ()-> {
			for (int i=1; i <= 1000; i++) {
				System.out.println("hello "+i);
			}
		};
		Runnable goodbyes = () -> {
			for (int i= 1; i <= 1000 ; i++) {
				System.out.println("goodbye "+i);
			}
		};
		ExecutorService executor = Executors.newCachedThreadPool();
		executor.execute(hellos);
		executor.execute(goodbyes);
	}
}
