package tacos;

import javax.persistence.Entity;
import javax.persistence.Id;

import lombok.Data;
import lombok.RequiredArgsConstructor;
import lombok.AccessLevel;
import lombok.NoArgsConstructor;

// https://projectlombok.org/
// java -jar lombok.jar
// eclipse.ini

@Data
@RequiredArgsConstructor
@NoArgsConstructor(access=AccessLevel.PRIVATE, force=true)
@Entity
public class Ingredient {

	// modifier final
	@Id
	private final String id;
    private final String name;
    private final Type type;
    public static enum Type {
      WRAP, PROTEIN, VEGGIES, CHEESE, SAUCE
}
}
