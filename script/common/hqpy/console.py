
import sys,platform
import logging

from colorama import init as coinit
from colorama import Fore, Back, Style

coinit(autoreset=True)
colormp = {logging.INFO:Fore.GREEN, logging.WARNING:Fore.YELLOW, logging.ERROR:Fore.RED}

# https://gist.github.com/curzona/9435729
class ProgressConsoleHandler(logging.StreamHandler):
    """
    A handler class which allows the cursor to stay on
    one line for selected messages
    """
    on_same_line = False
    linereturn = '\n'
    linefeed = '\r'
    indent = 0
    def emit(self, record):
        if hasattr(record, 'indent'):
            indent = getattr(record, 'indent')
            if indent:
                self.indent += 1
            else:
                self.indent -= 1
            return
                    
        try:
            msg = self.format(record)
            stream = self.stream
            same_line = hasattr(record, 'same_line')
            
            if self.on_same_line and not same_line:
                stream.write(self.linereturn)
            else:
                stream.write(self.linefeed)
                
            stream.write('  ' * self.indent + msg)
            
            if same_line:
                #stream.write('... ')
                self.on_same_line = True
            else:
                stream.write(self.linereturn)
                self.on_same_line = False
            
            self.flush()
        except (KeyboardInterrupt, SystemExit):
            raise
        except:
            #self.handleError(record)
            raise

class ColoramaConsoleHandler(logging.StreamHandler):
    btask_mode = False

    def setTaskMode(self):
        self.btask_mode = True

    def emit(self, record):       
        try:
            msg = self.format(record)
            if self.btask_mode == False and colormp.has_key(record.levelno) == True:
                cpre = colormp[record.levelno]
                msg = cpre  + msg + Fore.RESET
            self.stream.write(msg+'\n')
            self.flush()
        except (KeyboardInterrupt, SystemExit):
            raise
        except:
            #self.handleError(record)
            raise

def red(fmt, *arg):
    print(Fore.RED + fmt % (arg))

def green(fmt, *arg):
    print(Fore.GREEN + fmt % (arg))

def yellow(fmt, *arg):
    print(Fore.YELLOW + fmt % (arg))

def blue(fmt, *arg):
    print(Fore.BLUE + fmt % (arg))

def magenta(fmt, *arg):
    print(Fore.MAGENTA + fmt % (arg))

def cyan(fmt, *arg):
    print(Fore.CYAN + fmt % (arg))


