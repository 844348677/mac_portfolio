package chapter10;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class RaceCondition {
	public static void main(String[] args) {
		ExecutorService executor = Executors.newCachedThreadPool();
		int count = 0;
		for(int i=1 ; i <= 100 ; i++) {
			int taskID = i;
			Runnable task = () -> {
				for (int k=1; k<=1000; k++){
					//count = count +1;
				}
				System.out.println(taskID + " : "+count);
			};
			executor.execute(task);
		}
	}
}
