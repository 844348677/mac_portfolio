package ch05;

import java.util.Arrays;
import java.util.Comparator;
import java.util.List;
import static java.util.stream.Collectors.toList;


public class Practice {

	public static void main(String[] args) {
		// TODO Auto-generated method stub
		Trader raoul = new Trader("Raoul", "Cambridge");
        Trader mario = new Trader("Mario","Milan");
        Trader alan = new Trader("Alan","Cambridge");
        Trader brian = new Trader("Brian","Cambridge");
        List<Transaction> transactions = Arrays.asList(
            new Transaction(brian, 2011, 300),
            new Transaction(raoul, 2012, 1000),
            new Transaction(raoul, 2011, 400),
            new Transaction(mario, 2012, 710),
            new Transaction(mario, 2012, 700),
            new Transaction(alan, 2012, 950)
        		);		
        
        //Finds all transacsstions in 2011 and sort by value (small to high)
        List<Transaction> tr2011 =
        	    transactions.stream()
        	    .filter(transaction -> transaction.getYear() == 2011)
        	    .sorted(Comparator.comparing(Transaction::getValue))
        	    .collect(toList());
        
        //What are all the unique cities where the traders work?
        List<String> cities =
        	    transactions.stream()
        	    .map(transaction -> transaction.getTrader().getCity())
        	    .distinct()
        	    .collect(toList());
        /*
        Set<String> cities =
        	    transactions.stream()
        	                .map(transaction -> transaction.getTrader().getCity())
        	                .collect(toSet());
        */
        
	}

}
