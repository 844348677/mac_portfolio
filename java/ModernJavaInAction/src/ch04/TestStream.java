package ch04;

import java.util.Arrays;
import java.util.List;
import static java.util.stream.Collectors.toList;
import static java.util.stream.Collectors.*;



public class TestStream {
	
	
	public static void main(String[] args) {
		
		List<Dish> menu = Arrays.asList(
			    new Dish("pork", false, 800, Dish.Type.MEAT),
			    new Dish("beef", false, 700, Dish.Type.MEAT),
			    new Dish("chicken", false, 400, Dish.Type.MEAT),
			    new Dish("french fries", true, 530, Dish.Type.OTHER),
			    new Dish("rice", true, 350, Dish.Type.OTHER),
			    new Dish("season fruit", true, 120, Dish.Type.OTHER),
			    new Dish("pizza", true, 550, Dish.Type.OTHER),
			    new Dish("prawns", false, 300, Dish.Type.FISH),
			    new Dish("salmon", false, 450, Dish.Type.FISH) );
		
		List<String> threeHighCaloricDishNames =
				  menu.stream()
				      .filter(dish -> dish.getCalories() > 300)
				      .map(Dish::getName)
				      .limit(3)
				      .collect(toList());
		System.out.println(threeHighCaloricDishNames);
		
		List<String> names = 
				menu.stream()
					.filter(dish ->{
						System.out.println("filtering:" + dish.getName());
						return dish.getCalories() > 300;
					})
					.map(dish -> {
						System.out.println("mapping:" + dish.getName());
						return dish.getName();
					})
					.limit(3)
					.collect(toList());
		
		List<Integer> numbers = Arrays.asList(1, 2, 1, 3, 3, 2, 4);
		numbers.stream()
		       .filter(i -> i % 2 == 0)
		.distinct()
		       .forEach(System.out::println);
		
		List<Dish> specialMenu = Arrays.asList(
	            new Dish("seasonal fruit", true, 120, Dish.Type.OTHER),
	            new Dish("prawns", false, 300, Dish.Type.FISH),
	            new Dish("rice", true, 350, Dish.Type.OTHER),
	            new Dish("chicken", false, 400, Dish.Type.MEAT),
	            new Dish("french fries", true, 530, Dish.Type.OTHER));
		
		List<Dish> slicedMenu1 
				= specialMenu.stream()
					.takeWhile(dish -> dish.getCalories() < 320)
					.collect(toList());
		
		List<Dish> slicedMenu2
	    = specialMenu.stream()
	                 .dropWhile(dish -> dish.getCalories() < 320)
	                 .collect(toList());
		
		List<Integer> dishNameLengths = menu.stream()
				.map(Dish::getName)
				.map(String::length)
				.collect(toList());
		

		List<String> words = Arrays.asList("Hello","World");
		List<String> uniqueCharacters =
				  words.stream()
				       .map(word -> word.split(""))
				       .flatMap(Arrays::stream)
				       .distinct()
				       .collect(toList());
		
		List<Integer> numbers1 = Arrays.asList(1, 2, 3);
		List<Integer> numbers2 = Arrays.asList(3, 4);
		List<int[]> pairs =
		    numbers1.stream()
		            .flatMap(i -> numbers2.stream()
		                                  .map(j -> new int[]{i, j})
		)
		            .collect(toList());
		
		if(menu.stream().anyMatch(Dish::isVegetarian)) {
            System.out.println("The menu is (somewhat) vegetarian friendly!!");
		}
		

        boolean isHealthy = menu.stream()
                                .allMatch(dish -> dish.getCalories() < 1000);
        

        long count = menu.stream().count();
	}
}
