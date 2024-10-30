import { Link } from 'react-router-dom';
import Latex from 'react-latex';
import 'katex/dist/katex.min.css';
import {
  Container,
  Typography,
  Box,
  Paper,
  Stack,
  Divider,
  useTheme
} from '@mui/material';

function FaqPage() {
  const theme = useTheme();
  const brierScoreString = 'BS=\\frac{1}{N} \\sum^N_{t=1} (f_t-o_t)^2';
  const logScoreString = 'L({\\bf{r}}, i)=\\log_b(r_i)';
  const logScoreString2 = 'L({\\bf{r}}, i)=x\\log_b(p) + (1-x)\\log_b(1-p)';

  // Common styles for FAQ sections
  const sectionStyle = {
    p: 3,
    mb: 3,
    '& .katex': {
      fontSize: '1.2em',
    }
  };

  return (
    <Container 
      maxWidth="lg" 
      sx={{ 
        mt: { xs: 8, sm: 10 }, // Prevent header overlap
        mb: 4 
      }}
    >
      <Stack spacing={3}>
        <Paper sx={sectionStyle}>
          <Typography variant="h5" gutterBottom>FAQ 1: Why are you doing this?</Typography>
          <Typography paragraph>
            I want to make my beliefs match reality as much as possible. To make this happen, this public forecasting project is 
            a vital part. Now, my track record is visible to anyone—giving me an incentive to perform at my best in my forecasts.
          </Typography>
          <Typography paragraph>
            Eventually, I would like to add comments and the ability for readers to challenge my beliefs. Until then, you can 
            reach out to me with a question or challenge.
          </Typography>
        </Paper>

        <Paper sx={sectionStyle}>
          <Typography variant="h5" gutterBottom>FAQ 2: How are the scores you use defined?</Typography>
          <Typography paragraph>
            I use three different scoring rules: (i) Brier score, (ii) Binary Log score, and (iii) Natural Log score.
            They are defined as follows.
          </Typography>
          <Typography paragraph>
            (i) Brier score, when applied to one-dimensional forecasts, is the mean squared error of a forecast. This score is 
            minimized, so a score closer to 0 is better. A maximally uncertain prediction of 50% would yield a Brier score of 0.25.
          </Typography>
          <Box sx={{ my: 2, textAlign: 'center' }}>
            <Latex>{`$${brierScoreString}$`}</Latex>
          </Box>
          <Typography paragraph>
            (ii) Binary Log Score and (iii) Natural Log score are equivalent to the negative of Surprisal, a 
            commonly used scoring criterion in Bayesian inference. Compared to Brier score, these scores are harsher on poor
            predictions that are closer to 0 and 1. These scores are maximized, so a prediction scoring -0.10 is better than a 
            prediction scoring -0.8.
          </Typography>
          <Box sx={{ my: 2, textAlign: 'center' }}>
            <Latex>{`$${logScoreString}$`}</Latex>
          </Box>
          <Typography paragraph>
            This is a local strictly proper scoring rule under any base larger than 1. Given a binary prediction, the logarithmic 
            scoring rules can be rewritten as:
          </Typography>
          <Box sx={{ my: 2, textAlign: 'center' }}>
            <Latex>{`$${logScoreString2}$`}</Latex>
          </Box>
          <Typography>
            Where x is the outcome of the event and p is the expressed probability.
          </Typography>
        </Paper>

        <Paper sx={sectionStyle}>
          <Typography variant="h5" gutterBottom>FAQ 3: What is a proper scoring rule?</Typography>
          <Typography paragraph>
            A proper scoring rule is rule which, when implemented, incentivizes forecasters to report their True beliefs when forecasting.
            Mathematically, this means that a forecaster which maximizes reward will report the true probability distribution when they know the true
            probability distribution.
          </Typography>
        </Paper>

        <Paper sx={sectionStyle}>
          <Typography variant="h5" gutterBottom>FAQ 4: How can I contact you?</Typography>
          <Typography>
            You can reach me on email: bayesmaxxing@gmail.com
          </Typography>
        </Paper>

        <Paper sx={sectionStyle}>
          <Typography variant="h5" gutterBottom>FAQ 5: What are your plans for this project?</Typography>
          <Typography paragraph>
            Apart from design improvements and small features, such as free text search. I'm planning some larger features:
          </Typography>
          <Typography paragraph>
            (1) Include statistical models with my forecasts. I'm planning for this to be a link to a Google colab notebook or similar.
          </Typography>
          <Typography paragraph>
            (2) Comment section and Suggestion box for Users.
          </Typography>
          <Typography paragraph>
            (3) Blog posts.
          </Typography>
        </Paper>

        <Paper sx={sectionStyle}>
          <Typography variant="h5" gutterBottom>FAQ 6: Can I also forecast?</Typography>
          <Typography>
            Of course, anyone can do forecasting.
          </Typography>
        </Paper>

        <Paper sx={sectionStyle}>
          <Typography variant="h5" gutterBottom>FAQ 7: Why are many of the questions and predictions made on 2023-10-16?</Typography>
          <Typography>
            Before starting the project, I did not track the timing of my forecasts. To have a time of creation for all questions and predictions, 
            I set the time to the date I migrated my forecasts to the new database.
          </Typography>
        </Paper>
      </Stack>
    </Container>
  );
}

export default FaqPage;
