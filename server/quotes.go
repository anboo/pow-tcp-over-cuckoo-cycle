package main

import (
	"math/rand"
	"time"
)

var (
	wordOfWisdomQuotes = []string{
		"That inasmuch as any man drinketh wine or strong drink among you, behold it is not good, neither meet in the sight of your Father, only in assembling yourselves together to offer up your sacraments before him. (Doctrine and Covenants 89:5)",
		"And, behold, this should be wine, yea, pure wine of the grape of the vine, of your own make. (Doctrine and Covenants 89:6)",
		"And, again, strong drinks are not for the belly, but for the washing of your bodies. (Doctrine and Covenants 89:7)",
		"And again, tobacco is not for the body, neither for the belly, and is not good for man, but is an herb for bruises and all sick cattle, to be used with judgment and skill. (Doctrine and Covenants 89:8)",
		"And again, hot drinks are not for the body or belly. (Doctrine and Covenants 89:9)",
		"And again, verily I say unto you, all wholesome herbs God hath ordained for the constitution, nature, and use of man. (Doctrine and Covenants 89:10)",
		"Every herb in the season thereof, and every fruit in the season thereof; all these to be used with prudence and thanksgiving. (Doctrine and Covenants 89:11)",
		"Yea, flesh also of beasts and of the fowls of the air, I, the Lord, have ordained for the use of man with thanksgiving; nevertheless they are to be used sparingly. (Doctrine and Covenants 89:12)",
		"And it is pleasing unto me that they should not be used, only in times of winter, or of cold, or famine. (Doctrine and Covenants 89:13)",
		"All grain is ordained for the use of man and of beasts, to be the staff of life, not only for man but for the beasts of the field, and the fowls of heaven, and all wild animals that run or creep on the earth. (Doctrine and Covenants 89:14)",
		"All grain is good for the food of man; as also the fruit of the vine; that which yieldeth fruit, whether in the ground or above the ground. (Doctrine and Covenants 89:16)",
		"Nevertheless, wheat for man, and corn for the ox, and oats for the horse, and rye for the fowls and for swine, and for all beasts of the field, and barley for all useful animals, and for mild drinks, as also other grain. (Doctrine and Covenants 89:17)",
		"And all saints who remember to keep and do these sayings, walking in obedience to the commandments, shall receive health in their navel and marrow to their bones. (Doctrine and Covenants 89:18)",
		"And shall find wisdom and great treasures of knowledge, even hidden treasures. (Doctrine and Covenants 89:19)",
		"And shall run and not be weary, and shall walk and not faint. (Doctrine and Covenants 89:20)",
		"And I, the Lord, give unto them a promise, that the destroying angel shall pass by them, as the children of Israel, and not slay them. Amen. (Doctrine and Covenants 89:21)",
		"And all saints who remember to keep and do these sayings, walking in obedience to the commandments, shall receive health in their navel and marrow to their bones. (Doctrine and Covenants 89:18)",
		"Behold, verily, thus saith the Lord unto you: In consequence of evils and designs which do and will exist in the hearts of conspiring men in the last days, I have warned you, and forewarn you, by giving unto you this word of wisdom by revelation. (Doctrine and Covenants 89:4)",
		"That inasmuch as any man drinketh wine or strong drink among you, behold it is not good, neither meet in the sight of your Father, only in assembling yourselves together to offer up your sacraments before him. (Doctrine and Covenants 89:5)",
		"And, behold, this should be wine, yea, pure wine of the grape of the vine, of your own make. (Doctrine and Covenants 89:6)",
		"And, again, strong drinks are not for the belly, but for the washing of your bodies. (Doctrine and Covenants 89:7)",
		"And again, tobacco is not for the body, neither for the belly, and is not good for man, but is an herb for bruises and all sick cattle, to be used with judgment and skill. (Doctrine and Covenants 89:8)",
		"And again, hot drinks are not for the body or belly. (Doctrine and Covenants 89:9)",
		"And again, verily I say unto you, all wholesome herbs God hath ordained for the constitution, nature, and use of man. (Doctrine and Covenants 89:10)",
		"Every herb in the season thereof, and every fruit in the season thereof; all these to be used with prudence and thanksgiving. (Doctrine and Covenants 89:11)",
		"Yea, flesh also of beasts and of the fowls of the air, I, the Lord, have ordained for the use of man with thanksgiving; nevertheless they are to be used sparingly. (Doctrine and Covenants 89:12)",
		"And it is pleasing unto me that they should not be used, only in times of winter, or of cold, or famine. (Doctrine and Covenants 89:13)",
		"All grain is ordained for the use of man and of beasts, to be the staff of life, not only for man but for the beasts of the field, and the fowls of heaven, and all wild animals that run or creep on the earth. (Doctrine and Covenants 89:14)",
		"All grain is good for the food of man; as also the fruit of the vine; that which yieldeth fruit, whether in the ground or above the ground. (Doctrine and Covenants 89:16)",
		"Nevertheless, wheat for man, and corn for the ox, and oats for the horse, and rye for the fowls and for swine, and for all beasts of the field, and barley for all useful animals, and for mild drinks, as also other grain. (Doctrine and Covenants 89:17)",
		"And all saints who remember to keep and do these sayings, walking in obedience to the commandments, shall receive health in their navel and marrow to their bones. (Doctrine and Covenants 89:18)",
		"And shall find wisdom and great treasures of knowledge, even hidden treasures. (Doctrine and Covenants 89:19)",
		"And shall run and not be weary, and shall walk and not faint. (Doctrine and Covenants 89:20)",
		"And I, the Lord, give unto them a promise, that the destroying angel shall pass by them, as the children of Israel, and not slay them. Amen. (Doctrine and Covenants 89:21)",
		"And all saints who remember to keep and do these sayings, walking in obedience to the commandments, shall receive health in their navel and marrow to their bones. (Doctrine and Covenants 89:18)",
		"Behold, verily, thus saith the Lord unto you: In consequence of evils and designs which do and will exist in the hearts of conspiring men in the last days, I have warned you, and forewarn you, by giving unto you this word of wisdom by revelation. (Doctrine and Covenants 89:4)",
		"That inasmuch as any man drinketh wine or strong drink among you, behold it is not good, neither meet in the sight of your Father, only in assembling yourselves together to offer up your sacraments before him. (Doctrine and Covenants 89:5)",
		"And, behold, this should be wine, yea, pure wine of the grape of the vine, of your own make. (Doctrine and Covenants 89:6)",
		"And, again, strong drinks are not for the belly, but for the washing of your bodies. (Doctrine and Covenants 89:7)",
		"And again, tobacco is not for the body, neither for the belly, and is not good for man, but is an herb for bruises and all sick cattle, to be used with judgment and skill. (Doctrine and Covenants 89:8)",
		"And again, hot drinks are not for the body or belly. (Doctrine and Covenants 89:9)",
		"And again, verily I say unto you, all wholesome herbs God hath ordained for the constitution, nature, and use of man. (Doctrine and Covenants 89:10)",
		"Every herb in the season thereof, and every fruit in the season thereof; all these to be used with prudence and thanksgiving. (Doctrine and Covenants 89:11)",
		"Yea, flesh also of beasts and of the fowls of the air, I, the Lord, have ordained for the use of man with thanksgiving; nevertheless they are to be used sparingly. (Doctrine and Covenants 89:12)",
		"And it is pleasing unto me that they should not be used, only in times of winter, or of cold, or famine. (Doctrine and Covenants 89:13)",
	}

	randomizer = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func randomQuote() string {
	index := randomizer.Intn(len(wordOfWisdomQuotes))
	return wordOfWisdomQuotes[index]
}
