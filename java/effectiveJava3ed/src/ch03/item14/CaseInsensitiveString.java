package ch03.item14;

public class CaseInsensitiveString implements Comparable<CaseInsensitiveString> {
    private String s;
    @Override
    public int compareTo(CaseInsensitiveString cis) {
        return String.CASE_INSENSITIVE_ORDER.compare(s,cis.s);
    }
}
