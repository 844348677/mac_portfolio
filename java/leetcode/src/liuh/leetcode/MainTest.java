package liuh.leetcode;


import liuh.leetcode.summaryRanges228.Solution;

public class MainTest {
    public static void main(String[] args){
        /*
        771
        Solution jewels771 = new Solution();
        System.out.println(jewels771.numJewelsInStones("aA","aAAbbbb"));
        System.out.println(jewels771.numJewelsInStones("z","ZZ"));
        */

        Solution summaryRanges = new Solution();
        System.out.println(summaryRanges.summaryRanges(new int[]{0,1,2,4,5,7}));
        System.out.println(summaryRanges.summaryRanges(new int[]{0,2,3,4,6,8,9}));
    }
}
