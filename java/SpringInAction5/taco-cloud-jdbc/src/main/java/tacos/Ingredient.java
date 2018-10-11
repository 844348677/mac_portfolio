package tacos;

import lombok.Data;
import lombok.Getter;
import lombok.RequiredArgsConstructor;

// https://projectlombok.org/
// java -jar lombok.jar
// eclipse.ini

@Data
@RequiredArgsConstructor
public class Ingredient {

	// modifier final
	private final String id;
    private final String name;
    private final Type type;
    public static enum Type {
      WRAP, PROTEIN, VEGGIES, CHEESE, SAUCE
}
}
