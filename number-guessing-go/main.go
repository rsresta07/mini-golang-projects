package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	fmt.Println("--- WELCOME TO THE ROAST-O-MATIC GUESSING GAME ---")

	fmt.Println("\nChoose difficulty:")
	fmt.Println("1. Easy (1-50) 10 attempts")
	fmt.Println("2. Medium (1-100) 8 attempts")
	fmt.Println("3. Hard (1-500) 6 attempts")
	fmt.Println("4. Quit")

	for {
		var choice int
		fmt.Print("\nEnter choice (1-4): ")
		_, err := fmt.Scan(&choice)

		// Check if input was not a number
		if err != nil {
			fmt.Println("\n" + getInsult("invalid_type"))
			var discard string
			fmt.Scan(&discard) // Clear the bad input from buffer
			continue
		}

		var maxNumber, maxAttempts int
		switch choice {
		case 1:
			maxNumber = 50
			maxAttempts = 10
		case 2:
			maxNumber = 100
			maxAttempts = 8
		case 3:
			maxNumber = 500
			maxAttempts = 6
		case 4:
			fmt.Println("Giving up already? Typical. Goodbye!")
			return
		default:
			// If they picked a number like 99
			fmt.Println("\n" + getInsult("bad_menu_choice"))
			continue
		}

		secretNumber := rand.Intn(maxNumber) + 1

		fmt.Printf("\nI've picked a number between 1 and %d.\n", maxNumber)
		fmt.Printf("You have %d attempts. Good Luck (you'll need it).\n", maxAttempts)

		var guess int
		attemptsLeft := maxAttempts
		for attemptsLeft > 0 {
			fmt.Printf("\nAttempts left: %d -> Your guess: ", attemptsLeft)
			_, err := fmt.Scan(&guess)

			// If they typed letters during the game
			if err != nil {
				fmt.Println(getInsult("invalid_type"))
				continue
			}

			// If they guessed outside the range
			if guess < 1 || guess > maxNumber {
				fmt.Println(getInsult("out_of_range"))
				continue
			}

			attemptsLeft--

			if guess == secretNumber {
				fmt.Printf("\n%s\n", getInsult("win"))
				fmt.Printf("It took you %d attempts.\n", maxAttempts-attemptsLeft)
				break
			} else if guess < secretNumber {
				fmt.Println(getInsult("too_low"))
			} else {
				fmt.Println(getInsult("too_high"))
			}

			if attemptsLeft == 0 {
				// final burn
				fmt.Printf("\n%s\n", getInsult("game_over"))
				fmt.Printf("The number was %d by the way. Don't forget it.\n\n", secretNumber)
			}
		}

		fmt.Print("\nWanna lose again? (y/n): ")
		var again string
		fmt.Scan(&again)

		if again != "y" && again != "Y" && again != "yes" && again != "Yes" {
			fmt.Println("Finally, I can rest. See you next time!")
			break
		}
	}
}

// getInsult handles the emotional damage
func getInsult(category string) string {
	var list []string

	switch category {
	case "invalid_type": // When they type letters instead of numbers
		list = []string{
			"If confusion had a zip code, you just typed it out.",
			"That 'number' is so wrong it actually lowered the local GPS accuracy.",
			"Bruh, you didn't give me a number, you gave me a personality disorder.",
			"Bruh! That's not a number! Did you find those keys in a cereal box?",
			"I asked for a digit, not your life story in prose.",
			"Error 404: Brain not found. Please input an actual integer.",
			"Is your 'Enter' key the only thing working? Type a NUMBER.",
			"Your brain just tried to divide by zero, didn’t it?",
			"Quick question: In what dimension does that count as a value? I need to know so I can avoid it.",
			"That's not a number; that's a cry for help from your elementary school teacher.",
			"That's not a number, that's a symptom.",
			"Oh, cool. I didn't know we were doing creative writing in a math problem.",
			"That's not a digit, that's a typo of the soul.",
		}
	case "out_of_range": // When they guess 999 in a 1-50 game
		list = []string{
			"That number is so far out of bounds it's technically in a different time zone.",
			"Can you read? That's not within the range!",
			"Go back to school. That number is way out of bounds.",
			"You're guessing like someone who thinks the Earth is flat. Stay in the range!",
			"The range is right there on the screen. Are you blinking too much?",
			"Math is hard, but reading the instructions is free. Try it sometime!",
			"Is that your lucky number or just a random hallucination?",
			"You’re coloring outside the lines. Stay in the box, Picasso!",
			"Wait... do you actually know how to count, or are you just pressing buttons?",
			"That number is like your chances of winning: non-existent in this range.",
			"If 'clueless' was a coordinate, you'd be right where that number is.",
			"The boundaries aren't suggestions, they're the rules. Do you need a map?",
		}
	case "bad_menu_choice": // When they pick option 7 on a 1-4 menu
		list = []string{
			"There are 4 options. FOUR. How did you manage to mess that up?",
			"Selecting an invisible menu option? Bold strategy, let's see if it fails.",
			"That's not an option. Are you clicking buttons with your forehead?",
			"I gave you a simple 1 to 4 choice and you still choked.",
			"There are 4 options. 1, 2, 3, 4. Which part of that sequence is confusing you?",
			"That's not a menu choice, that's a cry for help.",
			"I gave you a map and you still walked into a wall. Pick 1, 2, 3, or 4!",
			"Is your keyboard broken, or are you just testing my patience?",
			"Congratulations! You found the 'Secret Option' that does absolutely nothing.",
			"If you handle your life like you handle this menu, I’m concerned for your future.",
			"I'm not saying you're bad at this, but my toaster is more logical than you.",
			"Error: User competence not found. Please select an ACTUAL option.",
			"Do you always just make things up as you go? Pick a number from the list!",
		}
	case "too_low":
		list = []string{
			"Too low. Are you trying to find the floor or just your self-esteem?",
			"Add some zeros, buddy. You’re currently in the basement.",
			"Too low! Just like your standards.",
			"Think bigger. Your guess is as tiny as your chances of winning.",
			"Wrong. Even your battery percentage is higher than that guess.",
			"Too low! You're playing like you're afraid of big numbers.",
			"If being wrong was an Olympic sport, you'd have the gold for that low-ball.",
			"Is that a guess or just the number of friends you have? Aim higher.",
			"That guess is so low it needs a submarine to find it.",
			"Aim higher! You're digging for rock bottom here.",
			"Too low. It's like you're not even trying.",
			"You're so close, yet still so disappointing.",
		}
	case "too_high":
		list = []string{
			"Too high! What are you sniffing?",
			"Too high! You’re not just out of range, you’re out of your mind.",
			"Calm down. That number is higher than the clouds you're living in.",
			"Is that a guess or your heart rate after trying to do basic math?",
			"That number is inflated. Just like your ego after getting one right.",
			"Too high! You're overshooting like a stormtrooper with bad eyesight.",
			"Too high. You're guessing like you're trying to win the lottery, not this game.",
			"Whoa there, NASA. Bring that guess back down to Earth.",
			"Bring it down a notch. You're in orbit right now.",
			"Too high! You're guessing like you have infinite attempts. (You don't).",
			"You're so close, yet still so disappointing.",
			"Calm down, Icarus. You're flying too close to the sun.",
		}
	case "game_over":
		list = []string{
			"Pathetic. My calculator could have done better.",
			"Game Over! Maybe try 'Connect 4'? This seems too hard for you.",
			"You lost. I'd say 'good game,' but I don't like lying.",
			"Mission failed. We'll get 'em next time... but probably not with you playing.",
			"Game Over. I’ve seen more intelligence in a 'Low Battery' warning.",
			"You failed. It’s a good thing this game doesn't require actual talent.",
			"Out of attempts? Stick to something easier, like blinking.",
			"You lost. I’d call you a 'noob,' but that feels like an insult to noobs everywhere.",
			"Game Over. You’re the reason they have to put instructions on shampoo bottles.",
			"You didn't find it. Don't worry, I didn't expect much from you anyway.",
			"Mission failed. I’d say 'better luck next time,' but luck isn't the problem here—competence is.",
			"You're out of tries. My CPU is literally yawning at your performance.",
		}
	case "win":
		list = []string{
			"AMAZING! You actually got it! Even a broken clock is right twice a day.",
			"You won? Wow, the bar was on the floor and you somehow cleared it.",
			"Correct! I assume you just guessed every number in your head until one stuck.",
			"You found it! I guess miracles really do happen for the technologically challenged.",
			"You actually found it. I guess if you throw enough darts at a board, eventually you'll hit the wall.",
			"Correct! I'll notify the local news that a miracle has occurred today.",
			"You won? I’m genuinely surprised. I already had the 'Loser' screen loading.",
			"Winner! It only took you long enough for the sun to go through a full cycle.",
			"You finally got it. I'd give you a trophy, but I don't think you could handle the weight.",
			"Correct. Even the most basic AI gets lucky once every trillion years.",
			"Wow, you found the number. Now try finding a clue.",
		}
	}

	return list[rand.Intn(len(list))]
}
