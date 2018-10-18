package ch06;

import java.util.Arrays;
import java.util.Comparator;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.Set;

import static java.util.stream.Collectors.*;
import static java.util.Arrays.asList;

import ch04.Dish;

public class Practice {

	public static void main(String[] args) {
		// TODO Auto-generated method stub

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

		Comparator<Dish> dishCaloriesComparator =
			    Comparator.comparingInt(Dish::getCalories);
			Optional<Dish> mostCalorieDish =
			    menu.stream()
			        .collect(maxBy(dishCaloriesComparator));
			
		int totalCalories = menu.stream().collect(summingInt(Dish::getCalories));
			

        String shortMenu = menu.stream().map(Dish::getName).collect(joining(", "));
        
        totalCalories = menu.stream().collect(reducing(0,
                Dish::getCalories,
                Integer::sum));
        
        Map<Dish.Type, List<Dish>> dishesByType =
                menu.stream().collect(groupingBy(Dish::getType));
        
        Map<CaloricLevel, List<Dish>> dishesByCaloricLevel = menu.stream().collect(
        groupingBy(dish -> {
        if (dish.getCalories() <= 400) return CaloricLevel.DIET;
        else if (dish.getCalories() <= 700) return CaloricLevel.NORMAL; else return CaloricLevel.FAT;
        } ));
        
        Map<String, List<String>> dishTags = new HashMap<>();
        dishTags.put("pork", asList("greasy", "salty"));
        dishTags.put("beef", asList("salty", "roasted"));
        dishTags.put("chicken", asList("fried", "crisp"));
        dishTags.put("french fries", asList("greasy", "fried"));
        dishTags.put("rice", asList("light", "natural"));
        dishTags.put("season fruit", asList("fresh", "natural"));
        dishTags.put("pizza", asList("tasty", "salty"));
        dishTags.put("prawns", asList("tasty", "roasted"));
        dishTags.put("salmon", asList("delicious", "fresh"));
        

		Map<Dish.Type, Set<String>> dishNamesByType =
		   menu.stream()
		      .collect(groupingBy(Dish::getType,
		               flatMapping(dish -> dishTags.get( dish.getName() ).stream(),
		            		   toSet())));
		Map<Dish.Type, Map<CaloricLevel, List<Dish>>> dishesByTypeCaloricLevel =
				menu.stream().collect( groupingBy(Dish::getType,
						groupingBy(dish -> {
							if (dish.getCalories() <= 400) return CaloricLevel.DIET;
							else if (dish.getCalories() <= 700) return CaloricLevel.NORMAL; else return CaloricLevel.FAT;
							}) )
				);
		
		Map<Dish.Type, Dish> mostCaloricByType =
			    menu.stream()
			    .collect(groupingBy(Dish::getType,
			            collectingAndThen(
			            		maxBy(Comparator.comparingInt(Dish::getCalories)),
			            		Optional::get)));
		
		Map<Boolean, List<Dish>> partitionedMenu =
	             menu.stream().collect(partitioningBy(Dish::isVegetarian));
		
		
			    
	}
	
	public enum CaloricLevel { DIET, NORMAL, FAT };

}
