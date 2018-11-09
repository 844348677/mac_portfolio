package liuh.leetcode.no228SummaryRanges;

import java.util.ArrayList;
import java.util.List;

public class Solution {
    public List<String> summaryRanges(int[] nums) {
        List<String> resultList = new ArrayList<>();
        if(nums.length==0)
            return resultList;
        if(nums.length==1) {
            resultList.add(nums[0] + "");
            return resultList;
        }
        boolean jump = false;
        int start = nums[0] ;
        for(int i=1;i<nums.length;i++){

            if(nums[i] > (nums[i-1]+1)){
                jump = true;
            }

            if(jump) {
                if (start == nums[i - 1])
                    resultList.add(start + "");
                else
                    resultList.add(start + "->" + nums[i-1]);
                start = nums[i];
                jump = false;
            }
        }
        if (start == nums[nums.length - 1])
            resultList.add(start + "");
        else
            resultList.add(start + "->" + nums[nums.length-1]);
        return resultList;
    }
}
/*
28.
Summary Ranges
Given a sorted integer array without duplicates, return the summary of its ranges.

Example 1:

Input:  [0,1,2,4,5,7]
Output: ["0->2","4->5","7"]
Explanation: 0,1,2 form a continuous range; 4,5 form a continuous range.
Example 2:

Input:  [0,2,3,4,6,8,9]
Output: ["0","2->4","6","8->9"]
Explanation: 2,3,4 form a continuous range; 8,9 form a continuous range.
 */
