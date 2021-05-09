package Saksuka

import "fmt"

func TokensToString(tokens []MessageToken) string {
	retval := "";
	fmt.Println(tokens)
	for _,v := range tokens {
		switch v.T {
			case "text": retval += v.V.(string);
			case "mention": retval += `@${`+v.V.(string)+`}`;
			case "link": retval += v.V.(string);
			case "emote": retval += `:${`+v.V.(string)+`}:`;
			case "block": retval +="`"+v.V.(string)+"`";
			default: retval += "";
		}
		retval += " ";
	}

	return retval;
}