package ch04.item15;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;

public class SecurityHole {
/*
Note that a nonzero-length array is always mutable, so it is wrong for a class to have a public static final array field, or an accessor that returns such a field.
 */
    // Potential security hole!
    public static final Object[] VALUES =  {};

    private static final Object[] PRIVATE_VALUES1 = { };
    public static final List<Object> VALUES1 =
        Collections.unmodifiableList(Arrays.asList(PRIVATE_VALUES1));

    private static final Object[] PRIVATE_VALUES2 = {};
    public static final Object[] values() {
        return PRIVATE_VALUES2.clone();
    }

}
