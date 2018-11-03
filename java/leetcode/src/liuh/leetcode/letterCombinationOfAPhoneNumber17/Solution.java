package liuh.leetcode.letterCombinationOfAPhoneNumber17;

import java.util.*;

public class Solution {

    public static final Map<String,List<String>> pair = new HashMap<String,List<String>>(){{
        put("2", Arrays.asList("a","b","c"));
        put("3",Arrays.asList("d","e","f"));
        put("4",Arrays.asList("g","h","i"));
        put("5",Arrays.asList("j","k","l"));
        put("6",Arrays.asList("m","n","o"));
        put("7",Arrays.asList("p","q","r","s"));
        put("8",Arrays.asList("t","u","v"));
        put("9",Arrays.asList("w","x","y","z"));
    }};

    public List<String> letterCombinations(String digits) {

        if(digits.equals(""))
            return new ArrayList<>();
        List<String> result = Arrays.asList(digits.split("")).parallelStream()
                .map(s -> {return pair.get(s);}).reduce(
                (first,second)->{
                    List<String> resultString = new ArrayList<>();
                    for(int i=0;i<first.size();i++){
                        for(int j=0;j<second.size();j++){
                            resultString.add(first.get(i)+second.get(j));
                        }
                    }
                    return resultString;
                }
        ).get();

        return result;
    }

}
/*
Given a string containing digits from 2-9 inclusive, return all possible letter combinations that the number could represent.

A mapping of digit to letters (just like on the telephone buttons) is given below. Note that 1 does not map to any letters.

 1 2 3
 4 5 6
 7 8 9

Input: "23"
Output: ["ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"].


 */
