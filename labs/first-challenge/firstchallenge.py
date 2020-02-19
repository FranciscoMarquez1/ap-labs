def longestSubstring (string):
    result = 0
    for i in range(len(string)):
        substring = set()
        counter = 0
        for j in range(i, len(string)):
            if string[j] not in substring:
                substring.add(string[j])
                counter += 1
            else:
                break
        if counter > result:
            result = counter
    return result

print(longestSubstring("pwwkew"))
