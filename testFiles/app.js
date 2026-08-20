const express = require('express');
const dbConfig = require('./dbConfig');

const app = express();
app.use(express.json());

app.get('/', async (req, res) => {
  try {
    const client = await dbConfig.pool.connect();
    const result = await client.query('SELECT NOW()');
    client.release();
    res.send(result.rows);
  } catch (err) {
    console.error(err.stack);
    res.status(500).send('Error connecting to database');
  }
});

app.listen(3000, () => {
  console.log('Server is running on port 3000');
});