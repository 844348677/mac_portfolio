package liuh.leetcode.no260SingleNumberIII;

/*
未完成
 */
public class Solution {
    public int[] singleNumber(int[] nums) {
        int resultsXOR = 0;
        for(int i=0;i<nums.length;i++){
            resultsXOR = resultsXOR ^ nums[i];
        }

        System.out.println(resultsXOR);


        return null;
    }
}
