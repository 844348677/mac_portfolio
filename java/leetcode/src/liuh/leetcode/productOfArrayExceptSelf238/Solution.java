package liuh.leetcode.productOfArrayExceptSelf238;

public class Solution {
    public int[] productExceptSelf(int[] nums) {

        int[] right = new int[nums.length];
        right[0] = 1;
        for(int i=1;i<nums.length;i++){
            right[i] = right[i-1] * nums[i-1];
        }
        int[] left = new int[nums.length];
        left[nums.length-1] = 1;
        for(int i= nums.length-2; i > -1 ;i--){
            left[i] = left[i+1] * nums[i+1];
        }
        int[] result = new int[nums.length];
        for(int i=0;i<nums.length;i++){
            result[i] = right[i] * left[i];
        }
        return result;
    }
}
/*
Given an array nums of n integers where n > 1,  return an array output such that output[i] is equal to the product of all the elements of nums except nums[i].

Example:

Input:  [1,2,3,4]
Output: [24,12,8,6]
Note: Please solve it without division and in O(n).

Follow up:
Could you solve it with constant space complexity? (The output array does not count as extra space for the purpose of space complexity analysis.)
 */
/*
        int[] result = new int[nums.length];

        for(int i=0;i<nums.length;i++){
            result[i] = 1;
            for(int j=0;j<nums.length;j++){
                if (i != j)
                    result[i] = result[i]*nums[j];
            }
        }
        return result;
 */