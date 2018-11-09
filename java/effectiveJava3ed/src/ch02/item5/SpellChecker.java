package ch02.item5;

import java.util.List;
import java.util.Objects;

class Lexicon{}

public class SpellChecker {
    private final Lexicon dictionary;
    public SpellChecker(Lexicon dictionary){
        this.dictionary = Objects.requireNonNull(dictionary);
    }
    public boolean isValid(String word){
        return false;
    }
    public List<String> suggestions(String typo){
        return null;
    }

}
