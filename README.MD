### GOLANG ORM for data management
### PYTHON FLASK for AI WORK and Automatic Messages
### MYSQL for database
### HTML-CSS-JS for frontend with "html/template" Golang package


### Task 1: Design the Review Chatbot

- **START:** Bot starts conversation by thanking user for buying and trusting the store, it addresses the user by name and ask nicely about the user's experience with the product by name (for example Iphone 13) on a manual message.
- **Messaging:** When the user answers, the bot will decide if the message talks about a refund, requesting information or rating the product or just something general, if its a refund or information request an automatic message will be sent if it's a general message an AI will answer and if it's a rating it will be saved to the tables throw gorm and a automatic response will be sent.

### Task 2: Model the Review Chatbot

- **Relational Database (MySQL) and GOLANG ORM (GORM):**

#### Tables:
- Users
- Conversations
- Orders
- OrderItems
- Reviews

#### Relationships:
- Users can have multiple Conversations and Orders.
- Conversations can have multiple Messages.
- Orders can have multiple OrderItems.
- OrderItems can have multiple Reviews.
- Reviews link back to both Users and OrderItems.

### Task 3: Implement the Review Chatbot

- **Demo Items: Includes three sample items (Phone, GPU, Monitor).**
- **User Flow: Users select an item and click "BUY," triggering the bot's initial message.**
- **Restriction: Access to "/" is restricted until an item is selected.**

### Task 4: Demonstrate the Chatbot

- **BUY Page: Initiates purchase and triggers the bot's conversation.**
- **Chatbot Page: Central hub for interactive sessions.**
- **User Page: Displays user-specific information**

---
