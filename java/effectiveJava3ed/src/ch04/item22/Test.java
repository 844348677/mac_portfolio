package ch04.item22;

// Use of static import to avoid qualifying constants
import static ch04.item22.PhysicalConstantsClass.AVOGADROS_NUMBER;

public class Test {
    double atoms(double mols) {
        return AVOGADROS_NUMBER * mols;
    }
    // Many more uses of PhysicalConstants justify static import }
}
