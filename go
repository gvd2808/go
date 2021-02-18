#define SPEED_L   6
#define SPEED_R   5 
#define DIR_L    7
#define DIR_R    4

// mooving speed
#define MSPEEDFAST 100
#define MSPEEDSLOW 25

// дальномер
#define PIN_TRIG 12
#define PIN_ECHO 13

#define LINELEFT  8
#define LINERIGHT  9

int blocked = 0;
long duration, dist;

void setup(){
  Serial.begin (9600);
  for( int i = 4; i != 10; i += 1){
    pinMode(i, OUTPUT);
  }
  pinMode(PIN_TRIG, OUTPUT);
  pinMode(PIN_ECHO, INPUT);
  
}

void mooveF(int Speed){
  analogWrite(SPEED_L, Speed+30);
  analogWrite(SPEED_R, Speed);

  digitalWrite(DIR_L, LOW);
  digitalWrite(DIR_R, HIGH);
}

void mooveB(int Speed){
  analogWrite(SPEED_L, Speed+30);
  analogWrite(SPEED_R, Speed);

  digitalWrite(DIR_L, HIGH);
  digitalWrite(DIR_R, LOW);
}

void turn_left(){
  analogWrite(SPEED_L,  100);
  analogWrite(SPEED_R, 100);

  digitalWrite(DIR_L, LOW);
  digitalWrite(DIR_R, LOW);
  delay(200);
  Mstop();
}

void turn_right(){
  analogWrite(SPEED_L, 100);
  analogWrite(SPEED_R, 100);

  digitalWrite(DIR_R, HIGH);
  digitalWrite(DIR_L, HIGH);
  delay(200);
  Mstop();
}
  
void Mstop(){
  // останавливаем робота
  analogWrite(SPEED_L, 0);
  analogWrite(SPEED_R, 0);
  Serial.println(blocked);
  digitalWrite(DIR_R, HIGH);
  digitalWrite(DIR_L, HIGH);
}

void find_dist(){
  long res = 0;
  // Сначала генерируем короткий импульс длительностью 2-5 микросекунд.
  for (int i = 0; i != 10; i++){
    digitalWrite(PIN_TRIG, LOW);
    delayMicroseconds(5);
    digitalWrite(PIN_TRIG, HIGH);
  
    // Выставив высокий уровень сигнала, ждем около 10 микросекунд. В этот момент датчик будет посылать сигналы с частотой 40 КГц.
    delayMicroseconds(10);
    digitalWrite(PIN_TRIG, LOW);
  
    //  Время задержки акустического сигнала на эхолокаторе.
    duration = pulseIn(PIN_ECHO, HIGH);
  
    // Теперь осталось преобразовать время в расстояние
    res = res + (duration / 58.2);
    delay(3);
  }
  
  dist = res / 10;
  Serial.print("dist is ");
  Serial.println(dist);
  
}

void tunnel_left(){
  while(dist > 20){
    find_dist();
    turn_left();
  }
  while(dist<20){
    find_dist();
    turn_left();
  }
  for(int i =0; i != 5; i++){
    turn_right();
  }
  Mstop();
  delay(5000);
}

void tunnel_right(){
  while(dist > 11){
    turn_right();
  }
  while(dist<11){
    turn_right();
  }
  for(int i =0; i != 5; i++){
    turn_right();
  }
  delay(5000);
}

void push(){
  while(digitalRead(LINELEFT)){
    mooveF(255);
  }
  delay(100);
  Mstop();
  find_dist();
  while(dist > 10){
    mooveB(100);
  }
}

void start_phase(){
  boolean state = true;
  find_dist();
  while(dist < 22){
    find_dist();
    if ( digitalRead(LINELEFT) == 0 and state == 1){
      ++blocked;
    }
    state = digitalRead(LINELEFT);
    mooveF(MSPEEDFAST);
  }
  Mstop();
}

void main_phase(){
  tunnel_left();
  for (int i = 1; i != 4; i++){
    if( i != blocked ){
      push();
      tunnel_right();
    }
  }
}

void finish(){
  while( digitalRead(LINELEFT)){
    mooveF(100);
  }
  Mstop();
}


void loop(){
  start_phase();
  main_phase();
  finish();
}
