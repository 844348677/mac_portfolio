package ch02.item2;

//  Calzone calzone = new Calzone.Builder()
//           .addTopping(HAM).sauceInside().build();

public class Calzone extends Pizza {
    private final boolean saceInside;

    public static class Builder extends Pizza.Builder<Builder> {
        private boolean sauceInside = false; //Default

        public Builder sauceInside(){
            sauceInside = true;
            return this;
        }

        @Override
        Pizza build() {
            return new Calzone(this);
        }

        @Override
        protected Builder self() {
            return this;
        }
    }
    Calzone(Builder builder) {
        super(builder);
        saceInside = builder.sauceInside;
    }
}
