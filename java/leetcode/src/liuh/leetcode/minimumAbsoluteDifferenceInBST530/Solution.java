package liuh.leetcode.minimumAbsoluteDifferenceInBST530;

import java.util.ArrayList;
import java.util.List;
import java.util.stream.Collectors;

/**
 * Definition for a binary tree node.
 * public class TreeNode {
 *     int val;
 *     TreeNode left;
 *     TreeNode right;
 *     TreeNode(int x) { val = x; }
 * }
 */

public class Solution {

    public int getMinimumDifference(TreeNode root) {

        List<Integer> allNodesVal = new ArrayList<>();
        List<Integer> tmp = binaryTreePreRecur(allNodesVal,root).parallelStream()
                .distinct().filter(data -> data != null).sorted().collect(Collectors.toList());
        int result=Math.abs(tmp.get(0)-tmp.get(1));
        for(int i=1;i<tmp.size();i++){
            int diff = tmp.get(i)-tmp.get(i-1);
            if(diff < result){
                 result = diff;
            }
        }

        return result;
    }

    public List<Integer> binaryTreePreRecur(List<Integer> allNodesVal, TreeNode root){
        if(root != null){
            allNodesVal.add(root.val);
            binaryTreePreRecur(allNodesVal,root.left);
            binaryTreePreRecur(allNodesVal,root.right);
        }
        return allNodesVal;
    }

}
/*
Given a binary search tree with non-negative values, find the minimum absolute difference between values of any two nodes.

Example:

Input:

   1
    \
     3
    /
   2

Output:
1

Explanation:
The minimum absolute difference is 1, which is the difference between 2 and 1 (or between 2 and 3).
 */