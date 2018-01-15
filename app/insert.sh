#!/bin/bash -e

if [ "$(curl -k https://localhost/v1/content)" == "null" ]
then
	curl -k https://localhost/v1/content -X POST -d "{
		\"title\":\"San Francisco: Something For Everyone\",
		\"body\":\"Home of tech behemoths some of which include AirBnB, Twitter, and Uber, San Francisco is perhaps most known for its booming industry. But what makes it such an iconic place to live? Whether it's the expansive lengths of the golden gate bridge, the classic display of Full House's very own painted ladies, or the towering presence of the newly constructed Salesforce tower - the likes of which can be seen from any point in the city - there's no arguing San Francisco is one of the most well known and beloved cities in America. Who knows what started the boom in San Francisco's popularity. Perhaps it was the top college in the area - U.C. Berkeley - that brought the best and brightest individuals to the city to study and to live. Or perhaps it was the location of the city within California; it's only a three hour drive from the city limits to the picturesque view of Yosemite National park, or an even shorter 16 miles to San Francisco's local stand of redwoods at Muir Woods National Monument. Plus, its location so close to the water means a temperate climate all year round. Whether it's the vibrant streets of the Castro, the outstanding culture, artwork, and architecture of the Mission District, or the bustling streets of the Financial district, San Francisco has a little something for everyone.\",
		\"extract\":\"Home of tech behemoths some of which include AirBnB, Twitter, and Uber, San Francisco is perhaps most known for its booming industry. But what makes it...\",
		\"location\":\"San Francisco, CA\",
		\"photo\":\"https://s3.amazonaws.com/mra-images/ngra-images/Default/golden-gate-bridge-1081782_1920.jpg\",
		\"author\":\"NGINX\"
	}"

	curl -k https://localhost/v1/content -X POST -d "{
		\"location\":\"Sydney, Australia\",
		\"author\":\"NGINX\",
		\"title\":\"Fun in Sydney\",
		\"extract\":\"Home to beaches, friendly people, and the world famous Opera House, Sydney is a must-see for any traveler...\",
		\"body\":\"The Australian cultural experience begins as soon as you set foot on the train between the airport and downtown. There you will meet locals and travelers alike.<br/>If you visit at the right time of the year, you can join in the festiviites which surround the Melbourne Cup, the biggest horse race in Australia. In Sydney, part of the pageantry of the event is for women to wear summer dresses with matching hats and for men to wear suits. If you find yourself in Sydney's Central Business District (or CBD) look for a bar called “Establishment”. There you will find a lot of friendly people who will show you a good time around Sydney.\",
		\"photo\":\"https://s3.amazonaws.com/mra-images/ngra-images/Default/sydney_harbor.jpg\"
	}"

	curl -k https://localhost/v1/content -X POST -d "{
		\"title\":\"Black and White Photos from Around the World\",
		\"body\":\"We have finally gotten around to putting up all the photos we have been taking recently and hopefully they will look better than ever. It's strange to be back from our month away - there's a reason why people say that the vacation blues are a thing. I never thought when my wife presented me with a one month trip to Europe that we would actually do it. No way would work let me off for such a lengthy amount of time, and of course the kids couldn't be left without their parents for so long. Looking back it all feels like a dream. We started in Paris, of course, being that was the place Shelly and I had always dreamed of visiting. Though it was a bit cold in the winter, that didn't stop us from walking the city streets at night and perusing the marketplaces places by day. We took our time in the city, spending about a week and a half wine tasting and dancing the nights away - but of course we didn't forget to relax as well. Then we went to the absolutely stunning mountains of Davos, Switzerland to ski. We booked a room in one of the hotels in the city, at the base of the mountain. We spent the next week delighting ourselves with scenic views on top of the mountain by day and cuddled up next to the fireplace built into our room at night - watching movies and drinking some of the best hot chocolate (secret ingredient? A little bit of grappa) I have ever tasted. It was a true bonding experience for the two of us; I really enjoy that we had this opportunity. We ended up coming back from the trip a week early because of homesickness. It was great to get away, but our true home is with the kids. Still, I'll never forget those nights in the city of love - streets lit by the festive decorations of the holidays. Or the eventful days atop the mountains of the Rhaetian Alps, skiing until our legs gave out from under us. Next we want to visit Vietnam for their delicious food and vibrant culture. But that's another day, another adventure.\",
		\"extract\":\"We have finally gotten around to putting up all the photos we have been taking recently and hopefully they will look better than ever. It's strange to be...\",
		\"location\":\"New York, NY\",
		\"photo\":\"https://s3.amazonaws.com/mra-images/ngra-images/Default/snow-2129837_1920.jpg\",
		\"author\":\"NGINX\"
	}"
fi
