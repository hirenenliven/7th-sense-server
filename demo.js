//routes.js
const express = require('express');
const router = express.Router();
const authController = require('../controllers/auth/authController');
const eventController = require('../controllers/event/eventController');
const { verifyToken } = require('../middleware/authMiddleware');
const { imapServiceMiddleware } = require('../middleware/serviceMiddleware');
const { validationMiddleware } = require('../middleware/validationMiddleware');

router.post('/login', validationMiddleware.login, authController.userLogin);
router.post('/logout', authController.logout);
router.post('/connection', validationMiddleware.verifyConnection, authController.verifyConnection);
router.post('/add-user', validationMiddleware.addUser, authController.addUser);
router.get('/userExist/:email', validationMiddleware.userExist, authController.userExist);
router.post('/verifyPassword', validationMiddleware.verifyPassword, authController.verifyPassword);
router.get('/verifyEmail/:email', validationMiddleware.userExist, authController.verifyEmail);
router.post('/verifyCode', validationMiddleware.verifyCode, authController.verifyCode);
router.get('/accept-decline-event', eventController.acceptDeclineEvent);

router.post('/sendVerificationCode', authController.sendVerificationCode);

module.exports = router;

//controller.js
const dotenv = require('dotenv');
const bcrypt = require('bcrypt');
const { verifyUserCredentials } = require('../../utils/account');
const connectionPool = require('../../utils/imapConnection');
const { verifyAppPasswordDB } = require('../../utils/account');
const { sendResponse } = require('../../common/sendResponse');
const { verifyEmailPasswordDB } = require('../../db/email/emails.service');
const { createUser, userExistDB, validateCredentials,verifyEmailDB,sendVerificationCodeEmail,verifyCodeDB } = require('../../db/user/user.service');
const { generateToken } = require('../../utils/commonFunctions');
const { createDefaultSettings } = require('../../db/settings/settings.service');
const { deleteUserSessionBySocketId } = require('../../db/userSessions/userSessions.service');

dotenv.config();

const userLogin = async (req, res) => {
    try {
        const userData = {
            role: req.body.role,
            email: req.body.email,
            password: req.body.password
        }

        const user = await validateCredentials(userData.email, userData.password, 'user');

        if (user === null) {
            return sendResponse(req, res, 401, { message: 'Invalid credentials' });
        }

        if (user.length === 0 || user.statusCode === 401 || user.statusCode === 404 || user.statusCode === 403) {
            return sendResponse(req, res, user.statusCode, {
                success: false,
                error: user.clientMessage || 'Authentication failed',
                statusCode: user.statusCode || 401
            });
        }

        return sendResponse(req, res, 200, {
            success: true,
            token: user.token,
            id: user.id,
            email: user.email,
            username: user.username
        });
    } catch (error) {
        return sendResponse(req, res, 500, { error: 'Failed to login' });
    }
}

//Response.js
const Responses = {
  success: {
    statusCode: 200,
    message: "success",
    clientMessage: { Message: "Success" },
  },
  created: {
    statusCode: 201,
    message: "created",
    clientMessage: { Message: "Created" },
  },
  emailExists: {
    statusCode: 400,
    message: "bad request",
    clientMessage: { Message: "Email already exists" },
  },
  adminExists: {
    statusCode: 400,
    message: "bad request",
    clientMessage: { Message: "Admin already exists" },
  },
  alreadyExist: {
    statusCode: 400,
    message: "bad request",
    clientMessage: { Message: "Already exists" },
  },
  badRequest: {
    statusCode: 400,
    message: "bad request",
    clientMessage: { Message: "Bad Request" },
  },
  unauthorized: {
    statusCode: 401,
    message: "unauthorized",
    clientMessage: { Message: "Unauthorized" },
  },
  paymentRequired: {
    statusCode: 402,
    message: "payment required",
    clientMessage: { Message: "Payment Required" },
  },
  notFound: {
    statusCode: 404,
    message: "not found",
    clientMessage: { Message: "Email not found" },
  },
  userNotFound: {
    statusCode: 404,
    message: "User not found",
    clientMessage: { Message: "User not found" },
  },
  serverTimeout: {
    statusCode: 408,
    message: "server timeout",
    clientMessage: { Message: "Server Timeout" },
  },
  requestEntryLarge: {
    statusCode: 413,
    message: "request entry too large",
    clientMessage: { Message: "Request entry too large" },
  },
  requestURLToLong: {
    statusCode: 414,
    message: "request-URL too long",
    clientMessage: { Message: "Request-URL too long" },
  },
  unsupportedMedia: {
    statusCode: 415,
    message: "unsupported media type",
    clientMessage: { Message: "Unsupported media type" },
  },
  expectationFailed: {
    statusCode: 417,
    message: "expectation failed",
    clientMessage: { Message: "Expectation failed" },
  },
  serverError: {
    statusCode: 500,
    message: "internal server error",
    clientMessage: { Message: "Internal server error" },
  },
  unavailable: {
    statusCode: 503,
    message: "service unavailable",
    clientMessage: { Message: "Service unavailable" },
  },
  validEmail: {
    statusCode: 400,
    message: "bad request",
    clientMessage: { Message: "Please enter valid email" },
  },
  validUserName: {
    statusCode: 400,
    message: "bad request",
    clientMessage: { Message: "Please enter valid UserName" },
  },
  validPassword: {
    statusCode: 400,
    message: "bad request",
    clientMessage: { Message: "Please enter valid password" },
  },
  tryAgain: {
    statusCode: 400,
    message: "bad request",
    clientMessage: { Message: "Something went wrong, Please try again later" },
  }
};

module.exports = { Responses };


//sendResponse.js
function sendJSONResponse(request, response, statusCode, messageOrObject) {
  let data;

  if (typeof messageOrObject === 'string') {
    data = { message: messageOrObject };
  } else if (typeof messageOrObject === 'object' && messageOrObject !== null) {
    data = messageOrObject;
  } else {
    data = { message: 'Invalid response format' };
  }

  response.status(statusCode).json({
    statusCode: statusCode,
    data: data
  });
}

function sendResponse(request, response, statusCode, messageOrObject) {
  sendJSONResponse(request, response, statusCode, messageOrObject);
}

module.exports = { sendResponse };


//user.service.js

const User = require('../../model/users.model');
const bcrypt = require('bcrypt');
const { generateToken } = require('../../utils/commonFunctions');
const { insertDocuments } = require('../../utils/dbUtils/dbCommonFunction');
const jwt = require('jsonwebtoken');
const { Responses } = require('../../common/responses');
const nodemailer = require('nodemailer');
const UserSession = require('../../model/userSessions.model');
const crypto = require('crypto');

async function createUser(userData) {
    try {
        // Check if user already exists
        const existingUser = await User.findOne({ email: userData.email });
        if (existingUser) {
            return Responses.emailExists;
        }

        if(userData.role === 'admin'){
            const isAlreadyAdmin = await User.findOne({ role: 'admin' });
            if (isAlreadyAdmin) {
                return Responses.adminExists;
            }
        }

        // Hash password
        const salt = await bcrypt.genSalt(10);
        const hashedPassword = await bcrypt.hash(userData.password, salt);

        userData.password = hashedPassword;

        const result = await User.create(userData);
        return result;
    } catch (error) {
        console.error('Error creating user:', error);
        return Responses.tryAgain;
    }
}