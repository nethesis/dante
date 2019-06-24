# Virgilio

## Configuration

### Environment variables

- VIRGILIO_ALLOW_ORIGIN
- VIRGILIO_STORE_DIR

## API endpoints

- `/miners`
  - `GET /` retrieve information on all the miners
- `/widget`
  - `GET /<widgetName>` retrieve information on the widget `<widgetName>`
- `layout`
  - `GET /` retrieve the layout for the widgets on the user dashboard. If not present, a default layout is returned
  - `POST /` save a custom layout for the widgets on the user dashboard
  - `DELETE /` delete the custom layout for the widgets on the user dashboard
- `/lang`
  - `GET /<langCode>` retrieve all the i18n strings for the language `<langCode>`
