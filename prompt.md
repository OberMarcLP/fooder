**Project Overview:**

I am building a web application with the following core features:

1. **Database**: PostgreSQL to store restaurant data.
2. **No Authentication**: The app will not require user authentication.
3. **Restaurant Rating System**: Users can rate restaurants based on food, service, and ambiance.
4. **Categories**: Support for cultural (e.g., Italian, Asian) and food-type (e.g., pizza, pasta) categories.
5. **Google Maps Integration**: Users can search for a restaurant by name, select the correct one from a list, and fetch details like name, address, and coordinates automatically from Google Maps. Additionally, users can view the restaurant's location on Google Maps.
6. **Edit/Delete Restaurants**: Users can edit or delete existing restaurants.
7. **Create/Edit/Delete Categories and Food Types**: Admins can create, edit, or delete cultural categories and food types.
8. **Theme Toggle**: A modern interface with support for both dark and light modes.

**Architecture & Technologies:**

* **Backend**: The backend will be built using a suitable programming language, suggested by you based on the use case.
* **Frontend**: React (modern UI with dark/light theme toggle).
* **Database**: PostgreSQL (storing restaurants, ratings, categories).
* **Google Maps API**: To search for restaurants by name, retrieve location details, and display the restaurant's location on the map.
* **Docker Compose**: To manage all services (backend, frontend, and database) in a containerized environment.

**Key Features:**

* **Restaurant Management**:

  * Users can search for a restaurant by name. The app will query the Google Maps API and display a list of matching restaurants.
  * Users will select the correct restaurant from the list. The app will fetch details such as name, address, coordinates, and description from Google Maps.
  * This information will automatically populate the restaurant form (name, description, address, food type, and category) for easy creation.
  * **Edit/Delete Restaurants**: Users can edit or delete existing restaurants.
* **Cultural Categories & Food Types**:

  * **Create/Edit/Delete Categories**: Admin users can create, edit, or delete cultural categories (e.g., Italian, Asian) and food types (e.g., pizza, pasta).
* **View Restaurant Location on Map**:

  * After selecting a restaurant, users can view the restaurant’s exact location on Google Maps embedded in the app.
  * Users can also view directions to the restaurant from their current location via Google Maps.
* **Rating System**: Users rate restaurants on food, service, and ambiance.
* **Map View**: Google Maps integration to show restaurant locations.
* **Theme Toggle**: A button to switch between dark and light modes.

**Google Maps Integration**:

* **Search Functionality**: Users can enter a restaurant's name into the search bar. The app will query the Google Maps API and display a list of matching restaurants.
* **Selection**: Users will select the correct restaurant from the list. The app will fetch details such as name, address, coordinates, and description from Google Maps.
* **View on Map**: After selecting a restaurant, the app will display a Google Map showing the exact location of the restaurant, with a marker showing its coordinates.
* **Directions**: Users can also get directions from their current location to the selected restaurant through Google Maps.

**Docker Setup:**

1. **Backend**: The backend will be containerized in Docker and will interact with PostgreSQL to handle restaurant data, ratings, Google Maps API calls, and CRUD operations for categories and food types.
2. **Frontend**: React app in a separate Docker container, displaying restaurant data and handling ratings and category management.
3. **PostgreSQL**: Runs in its own Docker container, storing restaurant and rating data.

